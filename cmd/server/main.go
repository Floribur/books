package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"

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
