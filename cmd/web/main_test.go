package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"snippetbox.mkz.net/internal/models"
)

func newTestApplication(t *testing.T) *application {
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	return &application{
		errorLog:      log.New(io.Discard, "", 0),
		infoLog:       log.New(io.Discard, "", 0),
		snippets:      models.NewSnippetModel(nil),
		templateCache: templateCache,
	}
}

func TestRoutes(t *testing.T) {
	app := newTestApplication(t)
	ts := httptest.NewServer(app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
	}{
		{
			name:     "Static asset exists",
			urlPath:  "/static/css/main.css",
			wantCode: http.StatusOK,
		},
		{
			name:     "Snippet create form",
			urlPath:  "/snippet/create",
			wantCode: http.StatusOK,
		},
		{
			name:     "Non-existent route returns 404",
			urlPath:  "/non-existent-route",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs, err := ts.Client().Get(ts.URL + tt.urlPath)
			if err != nil {
				t.Fatal(err)
			}
			defer rs.Body.Close()

			if rs.StatusCode != tt.wantCode {
				t.Errorf("got %d; want %d", rs.StatusCode, tt.wantCode)
			}
		})
	}
}

func TestServerSecureHeaders(t *testing.T) {
	app := newTestApplication(t)
	ts := httptest.NewServer(app.routes())
	defer ts.Close()

	rs, err := ts.Client().Get(ts.URL + "/static/css/main.css")
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	expectedHeaders := map[string]string{
		"Referrer-Policy":        "origin-when-cross-origin",
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "deny",
		"X-XSS-Protection":       "0",
	}

	for header, expectedValue := range expectedHeaders {
		if value := rs.Header.Get(header); value != expectedValue {
			t.Errorf("header %q got %q; want %q", header, value, expectedValue)
		}
	}
}

func TestServerStartupConfiguration(t *testing.T) {
	var buf bytes.Buffer
	errorLog := log.New(&buf, "ERROR\t", 0)

	// Verify openDB behaves safely when DSN is invalid / unavailable
	db, err := openDB("postgres://invalid:invalid@127.0.0.1:54399/snippetbox?sslmode=disable&connect_timeout=1")
	if err == nil {
		t.Fatal("expected error connecting to invalid DSN, got nil")
	}
	if db != nil {
		t.Errorf("expected db to be nil on failed connection, got %v", db)
	}
	errorLog.Printf("database connection warning: %v", err)
	if !bytes.Contains(buf.Bytes(), []byte("database connection warning")) {
		t.Errorf("expected log output to contain warning, got: %s", buf.String())
	}
}
