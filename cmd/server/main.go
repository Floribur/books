package main

import (
	"context"
	"fmt"
	"html"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"

	"flos-library/frontend"
	"flos-library/internal/api"
	"flos-library/internal/db"
	"flos-library/internal/scheduler"
	syncp "flos-library/internal/sync"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/floslib?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer pool.Close()

	queries := db.New(pool)

	// Embedded frontend filesystem
	distFS, err := fs.Sub(frontend.FS, "dist")
	if err != nil {
		log.Fatalf("failed to sub frontend dist: %v", err)
	}
	indexBytes, err := fs.ReadFile(distFS, "index.html")
	if err != nil {
		log.Fatalf("failed to read embedded index.html: %v", err)
	}
	indexHTML := string(indexBytes)

	// Enrichment trigger channel (buffered 1 — prevents double-trigger, RESEARCH.md Pitfall 7)
	enrichTrig := make(chan struct{}, 1)

	syncFn := func(ctx context.Context) error {
		err := syncp.SyncRSS(ctx, queries)
		if err != nil {
			log.Printf("sync error: %v", err)
			return err
		}
		// Trigger enrichment (non-blocking)
		select {
		case enrichTrig <- struct{}{}:
		default:
		}
		return nil
	}

	admin := &api.AdminHandlers{
		Queries:    queries,
		SyncFn:     syncFn,
		EnrichTrig: enrichTrig,
	}

	siteURL := os.Getenv("SITE_URL")
	if siteURL == "" {
		siteURL = "http://localhost:8081"
	}

	pub := &api.PublicHandlers{
		Store:   queries,
		SiteURL: siteURL,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	if os.Getenv("APP_ENV") == "development" {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins: []string{"http://localhost:5173"},
			AllowedMethods: []string{"GET", "OPTIONS"},
			AllowedHeaders: []string{"Accept", "Content-Type"},
			MaxAge:         300,
		}))
	}

	r.Post("/admin/sync", admin.PostSync)
	r.Post("/admin/import-csv", admin.PostImportCSV)

	// Book endpoints — order matters: specific before parameterized
	r.Get("/api/books/currently-reading", pub.GetCurrentlyReading)
	r.Get("/api/books", pub.GetBooks)
	r.Get("/api/books/{slug}", pub.GetBookBySlug)

	// Author endpoints
	r.Get("/api/authors", pub.GetAuthors)
	r.Get("/api/authors/{slug}", pub.GetAuthorBySlug)

	// Genre endpoints
	r.Get("/api/genres", pub.GetGenres)
	r.Get("/api/genres/{slug}", pub.GetGenreBySlug)

	// Years endpoint
	r.Get("/api/years", pub.GetYears)

	// Cover file server (API-09) — immutable cache headers for content-addressed images
	coversFS := http.Dir("data/covers")
	coversFileServer := http.FileServer(coversFS)
	r.Handle("/covers/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		http.StripPrefix("/covers/", coversFileServer).ServeHTTP(w, r)
	}))

	// SPA catch-all with OG meta tag injection (API-10)
	// /books/<slug> injects 5 Open Graph meta tags; all other non-API paths serve index.html.
	spaFileServer := http.FileServer(http.FS(distFS))
	r.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")

		// Serve real embedded files (JS, CSS, assets) directly without OG injection.
		if path != "" && path != "index.html" {
			f, err := distFS.Open(path)
			if err == nil {
				f.Close()
				spaFileServer.ServeHTTP(w, r)
				return
			}
		}

		// Detect /books/<slug> for OG tag injection (D-08).
		if strings.HasPrefix(r.URL.Path, "/books/") {
			slug := strings.Trim(strings.TrimPrefix(r.URL.Path, "/books/"), "/")
			if slug != "" {
				book, err := queries.GetBookDetailBySlug(r.Context(), slug)
				if err == nil {
					ogHTML := buildOGTags(book, siteURL)
					injected := strings.Replace(indexHTML, "</head>", ogHTML+"</head>", 1)
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					_, _ = w.Write([]byte(injected))
					return
				}
				// Book not found — fall through to serve plain index.html (SPA handles 404)
			}
		}

		// Default: serve index.html for all other routes (D-10).
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(indexHTML))
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go scheduler.Start(ctx, &wg, func(ctx context.Context) { _ = syncFn(ctx) })
	go syncp.RunEnricher(ctx, &wg, queries, enrichTrig)

	srv := &http.Server{Addr: ":" + port, Handler: r}
	go func() {
		log.Printf("server listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down...")
	_ = srv.Shutdown(context.Background())
	wg.Wait()
	log.Println("shutdown complete")
}

// buildOGTags constructs 5 Open Graph meta tags for a book detail page.
// Per decision D-08. Uses html.EscapeString to prevent XSS from book content.
// siteURL must not have a trailing slash.
func buildOGTags(book db.GetBookDetailBySlugRow, siteURL string) string {
	title := html.EscapeString(book.Title)

	desc := ""
	if book.Description != nil {
		d := *book.Description
		// Truncate at word boundary to ~200 chars (D-08 spec).
		if len(d) > 200 {
			truncated := d[:200]
			if idx := strings.LastIndex(truncated, " "); idx > 0 {
				truncated = truncated[:idx]
			}
			d = truncated + "\u2026" // "…"
		}
		desc = html.EscapeString(d)
	}

	ogImage := ""
	if book.CoverPath != nil && *book.CoverPath != "" {
		// cover_path stored as relative path like "data/covers/9781234567890.jpg"
		// Strip "data/covers/" prefix — covers are served at /covers/<filename>
		coverFile := strings.TrimPrefix(*book.CoverPath, "data/covers/")
		ogImage = fmt.Sprintf("%s/covers/%s", siteURL, coverFile)
	}

	pageURL := fmt.Sprintf("%s/books/%s", siteURL, book.Slug)

	return fmt.Sprintf(
		`<meta property="og:title" content="%s" />`+
			`<meta property="og:description" content="%s" />`+
			`<meta property="og:image" content="%s" />`+
			`<meta property="og:type" content="book" />`+
			`<meta property="og:url" content="%s" />`,
		title, desc, ogImage, pageURL,
	)
}
