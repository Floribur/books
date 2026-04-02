package scheduler

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestScheduler: mock ticker fires once → syncFn called exactly once after startup;
// context cancel → goroutine exits cleanly.
func TestScheduler(t *testing.T) {
	tickC := make(chan time.Time, 1)
	stopCalled := false
	stopFn := func() { stopCalled = true }

	var callCount atomic.Int32
	// syncDone is signaled after each syncFn call
	syncDone := make(chan struct{}, 10)
	syncFn := func(ctx context.Context) {
		callCount.Add(1)
		syncDone <- struct{}{}
	}

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)

	go startWithTicker(ctx, &wg, syncFn, tickC, stopFn)

	// Wait for the initial sync on startup
	select {
	case <-syncDone:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for initial sync")
	}
	if callCount.Load() != 1 {
		t.Errorf("expected 1 initial sync call, got %d", callCount.Load())
	}

	// Fire one tick
	tickC <- time.Now()

	// Wait for the tick-triggered sync
	select {
	case <-syncDone:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for tick sync")
	}
	if callCount.Load() != 2 {
		t.Errorf("expected 2 total calls after one tick, got %d", callCount.Load())
	}

	// Cancel context — goroutine should exit
	cancel()
	wg.Wait()

	if !stopCalled {
		t.Error("expected stop function to be called on exit")
	}
}
