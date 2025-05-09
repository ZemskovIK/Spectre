package sqlite

import (
	"database/sql"
	st "spectre/internal/storage"
)

const (
	GLOC = "src/internal/storage/sqlite/"
)

type sqliteDB struct {
	db *sql.DB
}

func New(dbPath string) (st.LettersStorage, error) {
	loc := GLOC + "New()"

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, ErrCannotConnectSQLite(loc, err)
	}

	return &sqliteDB{db: db}, nil
}

func (s *sqliteDB) Get(id int) (st.Letter, error) {
	return st.Letter{}, nil
}

func (s *sqliteDB) Save(letter st.Letter) error {
	return nil
}

func (s *sqliteDB) Delete(id int) error {
	return nil
}

func (s *sqliteDB) GetAll() ([]st.Letter, error) {
	return nil, nil
}
