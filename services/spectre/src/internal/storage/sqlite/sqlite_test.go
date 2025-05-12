package sqlite

import (
	"database/sql"
	"errors"
	"strings"
	"testing"

	st "spectre/internal/storage"
	"spectre/pkg/logger"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory SQLite database: %v", err)
	}

	// Create tables
	_, err = db.Exec(`
        CREATE TABLE authors (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            fname TEXT,
            mname TEXT,
            lname TEXT
        );
        CREATE TABLE letters (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT,
            body TEXT,
            found_at TEXT,
            found_in TEXT,
            author_id INTEGER,
            FOREIGN KEY(author_id) REFERENCES authors(id)
        );
    `)
	if err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}

	// Insert sample data
	_, err = db.Exec(`
        INSERT INTO authors (id, fname, mname, lname) VALUES (1, 'John', '', 'Doe');
        INSERT INTO letters (id, title, body, found_at, found_in, author_id)
        VALUES (1, 'Sample Title', 'Sample Body', '2023-01-01', 'Sample Location', 1);
    `)
	if err != nil {
		t.Fatalf("failed to insert sample data: %v", err)
	}

	return db
}

func TestGet_SuccessfullyRetrieveLetter(t *testing.T) {
	logger := logger.New("debug")
	db := setupTestDB(t)
	defer db.Close()

	storage := &sqliteDB{db: db, log: logger}

	letter, err := storage.Get(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if letter.ID != 1 || letter.Title != "Sample Title" || letter.Body != "Sample Body" {
		t.Errorf("unexpected letter: %+v", letter)
	}
}

func TestGet_LetterNotFound(t *testing.T) {
	logger := logger.New("debug")
	db := setupTestDB(t)
	defer db.Close()

	storage := &sqliteDB{db: db, log: logger}

	_, err := storage.Get(999)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	expectedErr := ErrLetterNotFound(999)
	if err.Error() != expectedErr.Error() {
		t.Errorf("expected error: %q, got: %q", expectedErr.Error(), err.Error())
	}
}

func TestGet_DatabaseError(t *testing.T) {
	logger := logger.New("debug")
	db := setupTestDB(t)
	defer db.Close()

	storage := &sqliteDB{db: db, log: logger}

	// Drop the table to simulate a database error
	_, err := storage.db.Exec("DROP TABLE letters")
	if err != nil {
		t.Fatalf("failed to drop table: %v", err)
	}

	_, err = storage.Get(1)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	// Check if the error is logged correctly
	if !errors.Is(err, sql.ErrNoRows) && !strings.Contains(err.Error(), "no such table") {
		t.Errorf("unexpected error: %v", err)
	}

	// Recreate the table to restore the database state
	_, err = storage.db.Exec(`
        CREATE TABLE letters (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT,
            body TEXT,
            found_at TEXT,
            found_in TEXT,
            author_id INTEGER,
            FOREIGN KEY(author_id) REFERENCES authors(id)
        );
    `)
	if err != nil {
		t.Fatalf("failed to recreate table: %v", err)
	}
}

func TestSave_SuccessfullySaveLetter(t *testing.T) {
	logger := logger.New("debug")
	db := setupTestDB(t)
	defer db.Close()

	storage := &sqliteDB{db: db, log: logger}

	letter := st.Letter{
		Title:   "New Title",
		Body:    "New Body",
		FoundAt: "2025-05-12",
		FoundIn: "New Location",
		Author:  "JaneDoe",
	}

	err := storage.Save(letter)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify the letter was saved
	var savedLetter st.Letter
	query := `SELECT l.id, l.title, l.body, l.found_at, l.found_in, a.fname || a.mname || a.lname AS author
              FROM letters l
              LEFT JOIN authors a ON l.author_id = a.id
              WHERE l.title = ? AND l.body = ?`
	err = db.QueryRow(query, letter.Title, letter.Body).Scan(
		&savedLetter.ID,
		&savedLetter.Title,
		&savedLetter.Body,
		&savedLetter.FoundAt,
		&savedLetter.FoundIn,
		&savedLetter.Author,
	)
	if err != nil {
		t.Fatalf("failed to retrieve saved letter: %v", err)
	}

	if savedLetter.Title != letter.Title || savedLetter.Body != letter.Body || savedLetter.Author != letter.Author {
		t.Errorf("saved letter does not match: %+v", savedLetter)
	}
}
