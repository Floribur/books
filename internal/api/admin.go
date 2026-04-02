package api

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"flos-library/internal/db"
	sync_pkg "flos-library/internal/sync"
)

// AdminHandlers holds dependencies for admin endpoints.
type AdminHandlers struct {
	Queries    *db.Queries
	SyncFn     func(context.Context) error
	EnrichTrig chan<- struct{}
	syncMu     sync.Mutex // prevents concurrent syncs
}

// PostSync handles POST /admin/sync
// Fires sync in a goroutine (fire-and-forget with 10min timeout), returns 202 immediately.
// Uses syncMu to prevent concurrent runs.
func (h *AdminHandlers) PostSync(w http.ResponseWriter, r *http.Request) {
	if !h.syncMu.TryLock() {
		http.Error(w, "sync already in progress", http.StatusConflict)
		return
	}
	go func() {
		defer h.syncMu.Unlock()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		if err := h.SyncFn(ctx); err != nil {
			// Log only — do not crash on manual trigger error
			_ = err
		}
		// Trigger enrichment after sync
		select {
		case h.EnrichTrig <- struct{}{}:
		default:
		}
	}()
	w.WriteHeader(http.StatusAccepted)
}

// PostImportCSV handles POST /admin/import-csv
// Accepts multipart/form-data with field "file".
func (h *AdminHandlers) PostImportCSV(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "invalid multipart form", http.StatusBadRequest)
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing file field", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Import runs synchronously so we can return an accurate count
	// CSV import for ~200 books completes in < 5 seconds
	count, err := sync_pkg.ImportCSV(r.Context(), h.Queries, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Trigger enrichment after import
	select {
	case h.EnrichTrig <- struct{}{}:
	default:
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]int{"imported": count})
}
