package models

import (
	"errors"
	"testing"
	"time"
)

func TestSnippetModelMemoryCRUD(t *testing.T) {
	m := NewSnippetModel(nil)

	// Test Latest contains seed data
	snippets, err := m.Latest()
	if err != nil {
		t.Fatalf("unexpected error getting latest: %v", err)
	}
	if len(snippets) != 2 {
		t.Fatalf("expected 2 seed snippets, got %d", len(snippets))
	}

	// Test Get existing
	s1, err := m.Get(1)
	if err != nil {
		t.Fatalf("unexpected error getting snippet 1: %v", err)
	}
	if s1.Title != "An old silent pond" {
		t.Errorf("expected title 'An old silent pond', got %q", s1.Title)
	}

	// Test Get non-existent
	_, err = m.Get(999)
	if !errors.Is(err, ErrNoRecord) {
		t.Errorf("expected ErrNoRecord for non-existent id, got %v", err)
	}

	// Test Insert
	id, err := m.Insert("New Title", "New Content", 7)
	if err != nil {
		t.Fatalf("unexpected error inserting snippet: %v", err)
	}
	if id != 3 {
		t.Errorf("expected id 3, got %d", id)
	}

	// Test Get newly inserted
	sNew, err := m.Get(id)
	if err != nil {
		t.Fatalf("unexpected error getting snippet %d: %v", id, err)
	}
	if sNew.Title != "New Title" || sNew.Content != "New Content" {
		t.Errorf("unexpected snippet content: %+v", sNew)
	}

	// Test Delete
	err = m.Delete(id)
	if err != nil {
		t.Fatalf("unexpected error deleting snippet: %v", err)
	}

	// Verify deleted
	_, err = m.Get(id)
	if !errors.Is(err, ErrNoRecord) {
		t.Errorf("expected ErrNoRecord after deletion, got %v", err)
	}

	// Delete non-existent
	err = m.Delete(999)
	if !errors.Is(err, ErrNoRecord) {
		t.Errorf("expected ErrNoRecord deleting non-existent id, got %v", err)
	}
}

func TestSnippetExpiration(t *testing.T) {
	m := &SnippetModel{
		memory: map[int]*Snippet{
			1: {
				ID:      1,
				Title:   "Expired Snippet",
				Content: "Content",
				Created: time.Now().UTC().Add(-48 * time.Hour),
				Expires: time.Now().UTC().Add(-24 * time.Hour),
			},
		},
		nextID: 2,
	}

	_, err := m.Get(1)
	if !errors.Is(err, ErrNoRecord) {
		t.Errorf("expected expired snippet to return ErrNoRecord, got %v", err)
	}

	snippets, err := m.Latest()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(snippets) != 0 {
		t.Errorf("expected 0 active snippets, got %d", len(snippets))
	}
}
