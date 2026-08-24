package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/shaardie/listinator/database"
)

// newTestServer builds a server backed by a fresh, migrated sqlite database
// in a temporary directory, following the same pattern as
// database.initTestDB in database/migration_test.go.
func newTestServer(t *testing.T) server {
	t.Helper()
	t.Setenv("LISTINATOR_ADMIN_PASSWORD", "")

	db, err := database.Init(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("unable to initialize database, %v", err)
	}
	return New(db)
}

// newTestRouter wires the real routes so requests exercise the exact same
// path/param binding as production, rather than hand-constructing an
// echo.Context.
func newTestRouter(s server) *echo.Echo {
	e := echo.New()
	s.SetupRoutes(e.Group(""))
	return e
}

// entryRequestBody covers every field the create/update/delete handlers
// accept; each test only sets the fields relevant to it. ID is a pointer so
// it's omitted from the JSON body for update/delete requests (uuid.UUID is
// a fixed-size array, so json's omitempty never treats a zero value as
// empty for it) — present, it would overwrite the path-bound ID for
// entryUpdate/entryDelete, whose input structs bind ID from both the URL
// param and (having no json tag of their own) the request body.
type entryRequestBody struct {
	ID              *uuid.UUID `json:"ID,omitempty"`
	Name            string     `json:"Name"`
	Number          string     `json:"Number"`
	Bought          bool       `json:"Bought"`
	TypeID          uuid.UUID  `json:"TypeID,omitempty"`
	ListID          uuid.UUID  `json:"ListID,omitempty"`
	ClientUpdatedAt *time.Time `json:"ClientUpdatedAt,omitempty"`
}

func doRequest(t *testing.T, e *echo.Echo, method, path string, body *entryRequestBody) (*httptest.ResponseRecorder, database.Entry) {
	t.Helper()

	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("unable to marshal request body, %v", err)
		}
		reader = bytes.NewReader(b)
	}

	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	var entry database.Entry
	if rec.Body.Len() > 0 {
		if err := json.Unmarshal(rec.Body.Bytes(), &entry); err != nil {
			t.Fatalf("unable to unmarshal response body %q, %v", rec.Body.String(), err)
		}
	}
	return rec, entry
}

func newEntry(t *testing.T, e *echo.Echo, clientUpdatedAt *time.Time) database.Entry {
	t.Helper()
	id := uuid.New()
	rec, entry := doRequest(t, e, http.MethodPost, "/entries", &entryRequestBody{
		ID:              &id,
		Name:            "first",
		TypeID:          uuid.New(),
		ListID:          uuid.New(),
		ClientUpdatedAt: clientUpdatedAt,
	})
	if rec.Code != http.StatusCreated {
		t.Fatalf("create: got status %v, want %v (body %v)", rec.Code, http.StatusCreated, rec.Body.String())
	}
	return entry
}

func TestEntryCreateIsIdempotent(t *testing.T) {
	s := newTestServer(t)
	e := newTestRouter(s)

	id := uuid.New()
	body := &entryRequestBody{
		ID:     &id,
		Name:   "milk",
		TypeID: uuid.New(),
		ListID: uuid.New(),
	}

	rec, first := doRequest(t, e, http.MethodPost, "/entries", body)
	if rec.Code != http.StatusCreated {
		t.Fatalf("first create: got status %v, want %v (body %v)", rec.Code, http.StatusCreated, rec.Body.String())
	}

	// Simulates the outbox replaying the same create after its original
	// response was lost.
	rec, replay := doRequest(t, e, http.MethodPost, "/entries", body)
	if rec.Code != http.StatusOK {
		t.Fatalf("replayed create: got status %v, want %v (body %v)", rec.Code, http.StatusOK, rec.Body.String())
	}
	if replay.ID != first.ID {
		t.Fatalf("replayed create returned a different row, got ID %v, want %v", replay.ID, first.ID)
	}

	var count int64
	if err := s.db.Model(&database.Entry{}).Where("id = ?", id).Count(&count).Error; err != nil {
		t.Fatalf("unable to count entries, %v", err)
	}
	if count != 1 {
		t.Fatalf("got %v rows for id %v, want 1", count, id)
	}
}

func TestEntryUpdateLastWriteWins(t *testing.T) {
	s := newTestServer(t)
	e := newTestRouter(s)

	// Anchored an hour in the past so every timestamp used below, including
	// the "second" edit, stays before the real "now" the nil-ClientUpdatedAt
	// subtest's request falls back to.
	t0 := time.Now().Add(-time.Hour)
	entry := newEntry(t, e, &t0)
	path := fmt.Sprintf("/entries/%v", entry.ID)

	t.Run("newer write is accepted", func(t *testing.T) {
		t1 := t0.Add(time.Minute)
		rec, updated := doRequest(t, e, http.MethodPut, path, &entryRequestBody{
			Name:            "second",
			TypeID:          entry.TypeID,
			ListID:          entry.ListID,
			ClientUpdatedAt: &t1,
		})
		if rec.Code != http.StatusOK {
			t.Fatalf("got status %v, want %v (body %v)", rec.Code, http.StatusOK, rec.Body.String())
		}
		if updated.Name != "second" {
			t.Fatalf("got name %q, want %q", updated.Name, "second")
		}
	})

	t.Run("stale write is silently rejected", func(t *testing.T) {
		tStale := t0.Add(30 * time.Second) // older than the "second" write above
		rec, rejected := doRequest(t, e, http.MethodPut, path, &entryRequestBody{
			Name:            "stale",
			TypeID:          entry.TypeID,
			ListID:          entry.ListID,
			ClientUpdatedAt: &tStale,
		})
		if rec.Code != http.StatusOK {
			t.Fatalf("got status %v, want %v (body %v)", rec.Code, http.StatusOK, rec.Body.String())
		}
		if rejected.Name != "second" {
			t.Fatalf("stale write was applied, got name %q, want unchanged %q", rejected.Name, "second")
		}
	})

	t.Run("nil ClientUpdatedAt always accepts", func(t *testing.T) {
		rec, updated := doRequest(t, e, http.MethodPut, path, &entryRequestBody{
			Name:   "no-timestamp",
			TypeID: entry.TypeID,
			ListID: entry.ListID,
		})
		if rec.Code != http.StatusOK {
			t.Fatalf("got status %v, want %v (body %v)", rec.Code, http.StatusOK, rec.Body.String())
		}
		if updated.Name != "no-timestamp" {
			t.Fatalf("got name %q, want %q", updated.Name, "no-timestamp")
		}
	})
}

func TestEntryDeleteConflict(t *testing.T) {
	s := newTestServer(t)
	e := newTestRouter(s)

	t.Run("delete at or after the stored timestamp succeeds", func(t *testing.T) {
		t0 := time.Now()
		entry := newEntry(t, e, &t0)
		path := fmt.Sprintf("/entries/%v", entry.ID)

		tDelete := t0.Add(time.Minute)
		rec, _ := doRequest(t, e, http.MethodDelete, path, &entryRequestBody{ClientUpdatedAt: &tDelete})
		if rec.Code != http.StatusOK {
			t.Fatalf("delete: got status %v, want %v (body %v)", rec.Code, http.StatusOK, rec.Body.String())
		}

		rec, _ = doRequest(t, e, http.MethodGet, path, nil)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("get after delete: got status %v, want %v", rec.Code, http.StatusNotFound)
		}
	})

	t.Run("stale delete is rejected and the newer edit survives", func(t *testing.T) {
		t0 := time.Now()
		entry := newEntry(t, e, &t0)
		path := fmt.Sprintf("/entries/%v", entry.ID)

		t1 := t0.Add(time.Minute)
		rec, _ := doRequest(t, e, http.MethodPut, path, &entryRequestBody{
			Name:            "second",
			TypeID:          entry.TypeID,
			ListID:          entry.ListID,
			ClientUpdatedAt: &t1,
		})
		if rec.Code != http.StatusOK {
			t.Fatalf("update: got status %v, want %v (body %v)", rec.Code, http.StatusOK, rec.Body.String())
		}

		tStaleDelete := t0.Add(30 * time.Second) // older than the "second" edit above
		rec, conflict := doRequest(t, e, http.MethodDelete, path, &entryRequestBody{ClientUpdatedAt: &tStaleDelete})
		if rec.Code != http.StatusConflict {
			t.Fatalf("stale delete: got status %v, want %v (body %v)", rec.Code, http.StatusConflict, rec.Body.String())
		}
		if conflict.Name != "second" {
			t.Fatalf("got name %q, want surviving edit %q", conflict.Name, "second")
		}

		rec, current := doRequest(t, e, http.MethodGet, path, nil)
		if rec.Code != http.StatusOK {
			t.Fatalf("get after rejected delete: got status %v, want %v", rec.Code, http.StatusOK)
		}
		if current.Name != "second" {
			t.Fatalf("entry did not survive the rejected delete, got name %q, want %q", current.Name, "second")
		}
	})

	t.Run("nil ClientUpdatedAt always proceeds", func(t *testing.T) {
		entry := newEntry(t, e, nil)
		path := fmt.Sprintf("/entries/%v", entry.ID)

		rec, _ := doRequest(t, e, http.MethodDelete, path, &entryRequestBody{})
		if rec.Code != http.StatusOK {
			t.Fatalf("delete: got status %v, want %v (body %v)", rec.Code, http.StatusOK, rec.Body.String())
		}

		rec, _ = doRequest(t, e, http.MethodGet, path, nil)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("get after delete: got status %v, want %v", rec.Code, http.StatusNotFound)
		}
	})
}
