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

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Post("/admin/sync", admin.PostSync)
	r.Post("/admin/import-csv", admin.PostImportCSV)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go scheduler.Start(ctx, &wg, func(ctx context.Context) { _ = syncFn(ctx) })

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
