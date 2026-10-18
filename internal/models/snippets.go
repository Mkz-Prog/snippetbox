package models

import (
	"database/sql"
	"errors"
	"sort"
	"sync"
	"time"
)

// Snippet represents a text snippet with title, content, and expiration.
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// SnippetModelInterface defines methods for snippet operations.
type SnippetModelInterface interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (*Snippet, error)
	Latest() ([]*Snippet, error)
	Delete(id int) error
}

// SnippetModel manages snippets in PostgreSQL with an in-memory fallback for local development.
type SnippetModel struct {
	DB     *sql.DB
	mu     sync.RWMutex
	memory map[int]*Snippet
	nextID int
}

// NewSnippetModel creates and initializes a new SnippetModel.
func NewSnippetModel(db *sql.DB) *SnippetModel {
	m := &SnippetModel{
		DB:     db,
		memory: make(map[int]*Snippet),
		nextID: 1,
	}

	if db == nil {
		now := time.Now().UTC()
		m.memory[1] = &Snippet{
			ID:      1,
			Title:   "An old silent pond",
			Content: "An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō",
			Created: now,
			Expires: now.Add(365 * 24 * time.Hour),
		}
		m.memory[2] = &Snippet{
			ID:      2,
			Title:   "Over the wintry forest",
			Content: "Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Sōseki",
			Created: now,
			Expires: now.Add(365 * 24 * time.Hour),
		}
		m.nextID = 3
	}

	return m
}

// Insert adds a new snippet to the database and returns its assigned ID.
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	if m.DB == nil {
		m.mu.Lock()
		defer m.mu.Unlock()

		id := m.nextID
		m.nextID++

		now := time.Now().UTC()
		s := &Snippet{
			ID:      id,
			Title:   title,
			Content: content,
			Created: now,
			Expires: now.Add(time.Duration(expires) * 24 * time.Hour),
		}
		if m.memory == nil {
			m.memory = make(map[int]*Snippet)
		}
		m.memory[id] = s
		return id, nil
	}

	stmt := `INSERT INTO snippets (title, content, created, expires)
    VALUES ($1, $2, (NOW() AT TIME ZONE 'UTC'), (NOW() AT TIME ZONE 'UTC') + ($3 * INTERVAL '1 day'))
    RETURNING id`

	var id int
	err := m.DB.QueryRow(stmt, title, content, expires).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Get retrieves a specific snippet by ID.
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	if m.DB == nil {
		m.mu.RLock()
		defer m.mu.RUnlock()

		s, exists := m.memory[id]
		if !exists || time.Now().UTC().After(s.Expires) {
			return nil, ErrNoRecord
		}
		copySnippet := *s
		return &copySnippet, nil
	}

	stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > (NOW() AT TIME ZONE 'UTC') AND id = $1`

	s := &Snippet{}
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return s, nil
}

// Latest returns the 10 most recently created unexpired snippets.
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	if m.DB == nil {
		m.mu.RLock()
		defer m.mu.RUnlock()

		now := time.Now().UTC()
		var snippets []*Snippet
		for _, s := range m.memory {
			if s.Expires.After(now) {
				copySnippet := *s
				snippets = append(snippets, &copySnippet)
			}
		}

		sort.Slice(snippets, func(i, j int) bool {
			return snippets[i].ID > snippets[j].ID
		})

		if len(snippets) > 10 {
			snippets = snippets[:10]
		}
		return snippets, nil
	}

	stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > (NOW() AT TIME ZONE 'UTC') ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}

		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

// Delete removes a snippet by ID.
func (m *SnippetModel) Delete(id int) error {
	if m.DB == nil {
		m.mu.Lock()
		defer m.mu.Unlock()

		if _, exists := m.memory[id]; !exists {
			return ErrNoRecord
		}
		delete(m.memory, id)
		return nil
	}

	stmt := `DELETE FROM snippets WHERE id = $1`
	result, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNoRecord
	}
	return nil
}
