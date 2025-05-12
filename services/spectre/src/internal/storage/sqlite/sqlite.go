package sqlite

import (
	"database/sql"
	st "spectre/internal/storage"
	"spectre/pkg/logger"
	"strings"
)

const (
	GLOC = "src/internal/storage/sqlite/"

	UNK_NAME = ""
)

type sqliteDB struct {
	db  *sql.DB
	log *logger.Logger
}

func New(dbPath string) (st.LettersStorage, error) {
	loc := GLOC + "New()"

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, ErrCannotConnectSQLite(loc, err)
	}

	return &sqliteDB{db: db}, nil
}

// Get retrieves a letter by id from db
func (s *sqliteDB) Get(id int) (st.Letter, error) {
	loc := GLOC + "Get()"

	query := `SELECT l.id, l.body, l.found_at, l.found_in, l.title
              FROM letters l
              LEFT JOIN authors a ON l.author_id = a.id
              WHERE l.id = ?`

	var letter st.Letter
	if err := s.db.QueryRow(query, id).Scan(
		&letter.ID,
		&letter.Body,
		&letter.FoundAt,
		&letter.FoundIn,
		&letter.Title,
	); err != nil {
		if err == sql.ErrNoRows {
			s.log.Warnf("%s: letter not found with id %d", loc, id)
			return st.Letter{}, ErrLetterNotFound(id)
		}
		s.log.Errorf("%s: error retrieving letter with id %d: %v", loc, id, err)
		return st.Letter{}, ErrCannotGetLetter(id, err)
	}

	s.log.Infof("%s: successfully retrieved letter with id %d", loc, id)
	return letter, nil
}

// Save saves a letter to db
func (s *sqliteDB) Save(letter st.Letter) error {
	loc := GLOC + "Save()"

	aID, err := s.getOrCreateAuthor(letter.Author)
	if err != nil {
		s.log.Errorf("%s: error getting or creating author '%s': %v", loc, letter.Author, err)
		return err
	}

	query := `INSERT INTO letters (title, body, found_at, found_in, author_id)
	          VALUES (?, ?, ?, ?, ?)`
	_, err = s.db.Exec(query,
		letter.Title,
		letter.Body,
		letter.FoundAt,
		letter.FoundIn,
		aID)
	if err != nil {
		s.log.Errorf("%s: error adding letter '%s': %v", loc, letter.Title, err)
		return ErrWhenAddingLetter(letter.Title, err)
	}

	s.log.Infof("%s: successfully saved letter '%s'", loc, letter.Title)
	return nil
}

func (s *sqliteDB) Delete(id int) error {
	loc := GLOC + "Delete()"
	s.log.Infof("%s: delete method called for id %d", loc, id)
	return nil
}

func (s *sqliteDB) GetAll() ([]st.Letter, error) {
	loc := GLOC + "GetAll()"
	s.log.Infof("%s: getAll method called", loc)
	return nil, nil
}

func splitName(name string) (string, string, string) {
	names := strings.Split(name, " ")
	switch len(names) {
	case 1:
		return names[0], UNK_NAME, UNK_NAME
	case 2:
		return names[0], names[1], UNK_NAME
	case 3:
		return names[0], names[1], names[2]
	default:
		return names[0], names[1], strings.Join(names[2:], " ")
	}
}

// getOrCreateAuthor checks if author exists in db, if not creates it
func (s *sqliteDB) getOrCreateAuthor(name string) (int, error) {
	loc := GLOC + "getOrCreateAuthor()"

	var fname, mname, lname string = splitName(name)
	var id int

	query := `SELECT id, fname, mname, lname
	          FROM authors
			  WHERE fname = ? AND mname = ? AND lname = ?`
	if err := s.db.QueryRow(query, fname, mname, lname).Scan(
		&id,
		&fname,
		&mname,
		&lname,
	); err == nil {
		s.log.Infof("%s: author '%s' found with id %d", loc, name, id)
		return id, nil
	} else if err != sql.ErrNoRows {
		s.log.Errorf("%s: error fetching author '%s': %v", loc, name, err)
		return -1, ErrWnenFetchAuthor(name, err)
	}

	query = `INSERT INTO authors (fname, mname, lname)
				VALUES (?, ?, ?)`
	res, err := s.db.Exec(query, fname, mname, lname)
	if err != nil {
		s.log.Errorf("%s: error adding author '%s': %v", loc, name, err)
		return -1, ErrWhenAddingAuthor(name, err)
	}

	uid, err := res.LastInsertId()
	if err != nil {
		s.log.Errorf("%s: error getting last insert id for author '%s': %v", loc, name, err)
		return -1, err
	}

	s.log.Infof("%s: successfully added author '%s' with id %d", loc, name, uid)
	return int(uid), nil
}
