package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func freshDb(t *testing.T, path ...string) *gorm.DB {
	t.Helper()

	var dbUri string

	// Note: path can be specified in an individual test for debugging
	// purposes -- so the db file can be inspected after the test runs.
	// Normally it should be left off so that a truly fresh memory db is
	// used every time.
	if len(path) == 0 {
		dbUri = ":memory:"
	} else {
		dbUri = path[0]
	}

	db, err := gorm.Open(sqlite.Open(dbUri), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error opening memory db: %s", err)
	}
	if err := setupDatabase(db); err != nil {
		t.Fatalf("Error setting up db: %s", err)
	}
	return db
}

// This tests that a fresh database returns no rows (but no error) when
// fetching Books.
func TestBookEmpty(t *testing.T) {
	db := freshDb(t)
	books := []Book{}
	if err := db.Find(&books).Error; err != nil {
		t.Fatalf("Error querying books from fresh db: %s", err)
	}
	if len(books) != 0 {
		t.Errorf("Expected 0 books, got %d", len(books))
	}
}

func TestDefaultRoute(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()

	ctx, r := gin.CreateTestContext(w)

	setupRouter(r, freshDb(t))

	req, err := http.NewRequestWithContext(ctx, "GET", "/", nil)
	if err != nil {
		t.Errorf("got error: %s", err)
	}

	r.ServeHTTP(w, req)
	if http.StatusOK != w.Code {
		t.Fatalf("expected response code %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()

	expected := "<h2>My Books</h2>"

	if !strings.Contains(body, expected) {
		t.Fatalf("expected response body '%s', got '%s'", expected, body)
	}
}
