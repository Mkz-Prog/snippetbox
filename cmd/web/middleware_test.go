package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeadersMiddleware(t *testing.T) {
	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	secureHeaders(next).ServeHTTP(rr, r)

	rs := rr.Result()
	defer rs.Body.Close()

	expectedHeaders := map[string]string{
		"Content-Security-Policy": "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com",
		"Referrer-Policy":         "origin-when-cross-origin",
		"X-Content-Type-Options":  "nosniff",
		"X-Frame-Options":         "deny",
		"X-XSS-Protection":        "0",
	}

	for header, expectedValue := range expectedHeaders {
		if value := rs.Header.Get(header); value != expectedValue {
			t.Errorf("header %q got %q; want %q", header, value, expectedValue)
		}
	}
}

func TestLogRequestMiddleware(t *testing.T) {
	var buf bytes.Buffer
	app := &application{
		infoLog: log.New(&buf, "", 0),
	}

	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/test-path", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.RemoteAddr = "127.0.0.1:12345"

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	app.logRequest(next).ServeHTTP(rr, r)

	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("HTTP/1.1 GET /test-path")) {
		t.Errorf("expected log output to contain request details, got: %s", output)
	}
}

func TestRecoverPanicMiddleware(t *testing.T) {
	app := &application{
		errorLog: log.New(io.Discard, "", 0),
	}

	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("something went horribly wrong")
	})

	app.recoverPanic(next).ServeHTTP(rr, r)

	rs := rr.Result()
	defer rs.Body.Close()

	if rs.StatusCode != http.StatusInternalServerError {
		t.Errorf("got status %d; want %d", rs.StatusCode, http.StatusInternalServerError)
	}

	if rs.Header.Get("Connection") != "close" {
		t.Errorf("expected Connection header 'close', got %q", rs.Header.Get("Connection"))
	}
}
