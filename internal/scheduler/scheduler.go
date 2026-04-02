package scheduler

import (
	"context"
	"log"
	"sync"
	"time"
)

const syncInterval = 6 * time.Hour

// Start runs syncFn immediately on startup, then every 6 hours.
// Exits cleanly when ctx is cancelled.
// Anti-pattern prevented: orphaned ticker via defer ticker.Stop().
func Start(ctx context.Context, wg *sync.WaitGroup, syncFn func(context.Context)) {
	ticker := time.NewTicker(syncInterval)
	startWithTicker(ctx, wg, syncFn, ticker.C, ticker.Stop)
}

// startWithTicker is the testable core of Start.
// It accepts an injectable tick channel and stop function.
func startWithTicker(ctx context.Context, wg *sync.WaitGroup, syncFn func(context.Context), tickC <-chan time.Time, stop func()) {
	defer wg.Done()
	defer stop()

	// Run immediately on startup
	log.Println("scheduler: initial sync on startup")
	syncFn(ctx)

	for {
		select {
		case <-tickC:
			log.Println("scheduler: tick, running sync")
			syncFn(ctx)
		case <-ctx.Done():
			log.Println("scheduler: context cancelled, exiting")
			return
		}
	}
}
