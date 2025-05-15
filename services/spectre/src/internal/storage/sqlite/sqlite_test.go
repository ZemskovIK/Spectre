package sqlite

import (
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	st "spectre/internal/storage"
	"spectre/pkg/logger"

	_ "github.com/mattn/go-sqlite3"
)

func mustParseDate(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return t
}

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
            fname TEXT NOT NULL,
            mname TEXT NOT NULL,
            lname TEXT NOT NULL
        );
        CREATE TABLE letters (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            body TEXT NOT NULL,
            found_at TIMESTAMP NOT NULL,
            found_in TEXT NOT NULL,
            author_id INTEGER DEFAULT NULL,
            FOREIGN KEY(author_id) REFERENCES authors(id) ON DELETE SET NULL
        );
    `)
	if err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}

	// Insert sample data
	_, err = db.Exec(`
        INSERT INTO authors (id, fname, mname, lname) VALUES (1, 'John', 'Doe', '');
        INSERT INTO letters (id, body, found_at, found_in, author_id)
        VALUES (1, 'Sample Body', '2023-01-01', 'Sample Location', 1);
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

	if letter.ID != 1 || letter.Body != "Sample Body" {
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

	expectedErr := errLetterNotFound(999)
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
            body TEXT NOT NULL,
            found_at TEXT NOT NULL,
            found_in TEXT NOT NULL,
            author_id INTEGER DEFAULT NULL,
            FOREIGN KEY(author_id) REFERENCES authors(id) ON DELETE SET NULL
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
		Body:    "New Body",
		FoundAt: time.Now(),
		FoundIn: "New Location",
		Author:  "Jane Doe",
	}

	err := storage.Save(letter)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify the letter was saved
	var savedLetter st.Letter
	query := `SELECT l.id, l.body, l.found_at, l.found_in, 
                     TRIM(a.fname || ' ' || a.mname || ' ' || a.lname) AS author
              FROM letters l
              LEFT JOIN authors a ON l.author_id = a.id
              WHERE l.body = ?`
	err = db.QueryRow(query, letter.Body).Scan(
		&savedLetter.ID,
		&savedLetter.Body,
		&savedLetter.FoundAt,
		&savedLetter.FoundIn,
		&savedLetter.Author,
	)
	if err != nil {
		t.Fatalf("failed to retrieve saved letter: %v", err)
	}

	if savedLetter.Body != letter.Body || savedLetter.Author != letter.Author {
		t.Errorf("saved letter does not match: %+v", savedLetter)
	}
}

func TestGetAll_SuccessfullyRetrieveAllLetters(t *testing.T) {
	logger := logger.New("debug")
	db := setupTestDB(t)
	defer db.Close()

	storage := &sqliteDB{db: db, log: logger}

	// Добавляем авторов с разными форматами имен
	_, err := db.Exec(`
        INSERT INTO authors (id, fname, mname, lname) 
        VALUES 
            (2, 'Alice', 'Smith', ''),
            (3, 'Bob', 'Michael', 'Asd'),
            (4, 'Charlie', '', '');

        INSERT INTO letters (id, body, found_at, found_in, author_id)
        VALUES 
            (2, 'Second Body', '2025-05-12', 'Second Location', 2),
            (3, 'Third Body', '2025-05-13', 'Third Location', 3),
            (4, 'Fourth Body', '2025-05-14', 'Fourth Location', 4);
    `)
	if err != nil {
		t.Fatalf("failed to insert additional authors and letters: %v", err)
	}

	// Вызываем метод GetAll
	letters, err := storage.GetAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Проверяем количество полученных писем
	if len(letters) != 4 {
		t.Fatalf("expected 4 letters, got %d", len(letters))
	}

	// Проверяем содержимое писем
	expectedLetters := []st.Letter{
		{
			ID:      1,
			Body:    "Sample Body",
			FoundAt: mustParseDate("2023-01-01"),
			FoundIn: "Sample Location",
			Author:  "John Doe",
		},
		{
			ID:      2,
			Body:    "Second Body",
			FoundAt: mustParseDate("2025-05-12"),
			FoundIn: "Second Location",
			Author:  "Alice Smith",
		},
		{
			ID:      3,
			Body:    "Third Body",
			FoundAt: mustParseDate("2025-05-13"),
			FoundIn: "Third Location",
			Author:  "Bob Michael Asd",
		},
		{
			ID:      4,
			Body:    "Fourth Body",
			FoundAt: mustParseDate("2025-05-14"),
			FoundIn: "Fourth Location",
			Author:  "Charlie",
		},
	}
	for i, expected := range expectedLetters {
		if letters[i] != expected {
			t.Errorf("expected letter at index %d: %+v, got: %+v", i, expected, letters[i])
		}
	}
}

func TestDelete(t *testing.T) {
	logger := logger.New("debug")
	db := setupTestDB(t)
	defer db.Close()

	storage := &sqliteDB{db: db, log: logger}

	t.Run("successfully delete an existing letter", func(t *testing.T) {
		// Удаляем существующее письмо
		err := storage.Delete(1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Проверяем, что письмо удалено
		_, err = storage.Get(1)
		if err.Error() != errLetterNotFound(1).Error() {
			t.Errorf("expected ErrLetterNotFound, got: %v", err)
		}
	})

	t.Run("delete a non-existing letter", func(t *testing.T) {
		// Пытаемся удалить несуществующее письмо
		err := storage.Delete(999)
		if err.Error() != errLetterNotFound(999).Error() {
			t.Errorf("expected ErrLetterNotFound, got: %v", err)
		}
	})
}
