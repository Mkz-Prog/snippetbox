package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHome(t *testing.T) {
	app := newTestApplication(t)
	ts := httptest.NewServer(app.routes())
	defer ts.Close()

	rs, err := ts.Client().Get(ts.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	if rs.StatusCode != http.StatusOK {
		t.Errorf("got status %d; want %d", rs.StatusCode, http.StatusOK)
	}

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(body), "Latest Snippets") {
		t.Errorf("expected body to contain 'Latest Snippets'")
	}
	if !strings.Contains(string(body), "An old silent pond") {
		t.Errorf("expected body to contain seed snippet 'An old silent pond'")
	}
}

func TestSnippetView(t *testing.T) {
	app := newTestApplication(t)
	ts := httptest.NewServer(app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "An old silent pond",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/snippet/view/999",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/snippet/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/snippet/view/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/snippet/view/abc",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID (route mismatch)",
			urlPath:  "/snippet/view/",
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
				t.Errorf("got status %d; want %d", rs.StatusCode, tt.wantCode)
			}

			if tt.wantBody != "" {
				body, err := io.ReadAll(rs.Body)
				if err != nil {
					t.Fatal(err)
				}
				if !strings.Contains(string(body), tt.wantBody) {
					t.Errorf("expected body to contain %q", tt.wantBody)
				}
			}
		})
	}
}

func TestSnippetCreate(t *testing.T) {
	app := newTestApplication(t)
	ts := httptest.NewServer(app.routes())
	defer ts.Close()

	rs, err := ts.Client().Get(ts.URL + "/snippet/create")
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	if rs.StatusCode != http.StatusOK {
		t.Errorf("got status %d; want %d", rs.StatusCode, http.StatusOK)
	}

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(body), "<form action='/snippet/create' method='POST'>") {
		t.Errorf("expected create snippet form in body")
	}
}

func TestSnippetCreatePost(t *testing.T) {
	app := newTestApplication(t)
	ts := httptest.NewServer(app.routes())
	defer ts.Close()

	// Disable redirect following so we can inspect 303 StatusSeeOther
	client := ts.Client()
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	tests := []struct {
		name         string
		title        string
		content      string
		expires      string
		wantCode     int
		wantLocation string
		wantBody     string
	}{
		{
			name:         "Valid submission",
			title:        "Automated Test Snippet",
			content:      "Testing snippet content here",
			expires:      "7",
			wantCode:     http.StatusSeeOther,
			wantLocation: "/snippet/view/3",
		},
		{
			name:     "Empty title",
			title:    "",
			content:  "Some content",
			expires:  "7",
			wantCode: http.StatusUnprocessableEntity,
			wantBody: "This field cannot be blank",
		},
		{
			name:     "Title too long",
			title:    strings.Repeat("a", 101),
			content:  "Some content",
			expires:  "7",
			wantCode: http.StatusUnprocessableEntity,
			wantBody: "This field cannot be more than 100 characters long",
		},
		{
			name:     "Empty content",
			title:    "Valid Title",
			content:  "",
			expires:  "7",
			wantCode: http.StatusUnprocessableEntity,
			wantBody: "This field cannot be blank",
		},
		{
			name:     "Invalid expires",
			title:    "Valid Title",
			content:  "Valid content",
			expires:  "99",
			wantCode: http.StatusUnprocessableEntity,
			wantBody: "This field must equal 1, 7 or 365",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("title", tt.title)
			form.Add("content", tt.content)
			form.Add("expires", tt.expires)

			rs, err := client.PostForm(ts.URL+"/snippet/create", form)
			if err != nil {
				t.Fatal(err)
			}
			defer rs.Body.Close()

			if rs.StatusCode != tt.wantCode {
				t.Errorf("got status %d; want %d", rs.StatusCode, tt.wantCode)
			}

			if tt.wantLocation != "" {
				location := rs.Header.Get("Location")
				if location != tt.wantLocation {
					t.Errorf("got Location %q; want %q", location, tt.wantLocation)
				}
			}

			if tt.wantBody != "" {
				body, err := io.ReadAll(rs.Body)
				if err != nil {
					t.Fatal(err)
				}
				if !strings.Contains(string(body), tt.wantBody) {
					t.Errorf("expected body to contain %q", tt.wantBody)
				}
			}
		})
	}
}
