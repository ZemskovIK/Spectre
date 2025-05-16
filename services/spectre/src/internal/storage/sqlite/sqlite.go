package sqlite

import (
	"database/sql"
	"spectre/internal/lib"
	st "spectre/internal/storage"
	"spectre/pkg/logger"

	_ "github.com/mattn/go-sqlite3"
)

const (
	GLOC = "src/internal/storage/sqlite/"

	UNK_NAME = ""
)

type sqliteDB struct {
	db  *sql.DB
	log *logger.Logger
}

func New(dbPath string, log *logger.Logger) (st.LettersStorage, error) {
	loc := GLOC + "New()"
	log.Debugf("%s: opening sqlite db at path: %s", loc, dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Errorf("%s: failed to open sqlite db: %v", loc, err)
		return nil, errCannotConnectSQLite(loc, err)
	}

	return &sqliteDB{
		db:  db,
		log: log,
	}, nil
}

// Get retrieves a letter by id from db
func (s *sqliteDB) Get(id int) (st.Letter, error) {
	loc := GLOC + "Get()"
	s.log.Debugf("%s: preparing to get letter with id: %d", loc, id)

	query := `SELECT l.id, l.body, l.found_at, l.found_in
              FROM letters l
              LEFT JOIN authors a ON l.author_id = a.id
              WHERE l.id = ?`
	s.log.Debugf("%s: executing query: %s with id=%d", loc, query, id)

	var letter st.Letter
	if err := s.db.QueryRow(query, id).Scan(
		&letter.ID,
		&letter.Body,
		&letter.FoundAt,
		&letter.FoundIn,
	); err != nil {
		if err == sql.ErrNoRows {
			s.log.Warnf("%s: letter not found with id %d", loc, id)
			return st.Letter{}, errLetterNotFound(id)
		}
		s.log.Errorf("%s: error retrieving letter with id %d: %v", loc, id, err)
		return st.Letter{}, errCannotGetLetter(id, err)
	}

	s.log.Infof("%s: successfully retrieved letter with id %d", loc, id)
	return letter, nil
}

// Save saves a letter to db
func (s *sqliteDB) Save(letter st.Letter) error {
	loc := GLOC + "Save()"
	s.log.Debugf("%s: preparing to save letter: %+v", loc, letter)

	aID, err := s.getOrCreateAuthor(letter.Author)
	if err != nil || aID <= 0 {
		s.log.Errorf("%s: error getting or creating author '%s': %v", loc, letter.Author, err)
		return err
	}
	s.log.Debugf("%s: using author_id=%d for letter", loc, aID)

	query := `INSERT INTO letters (body, found_at, found_in, author_id)
	          VALUES (?, ?, ?, ?)`
	s.log.Debugf("%s: executing query: %s with values: body=%s, found_at=%s, found_in=%s, author_id=%d", loc, query, letter.Body, letter.FoundAt, letter.FoundIn, aID)
	_, err = s.db.Exec(query,
		letter.Body,
		letter.FoundAt,
		letter.FoundIn,
		aID)
	if err != nil {
		s.log.Errorf("%s: error adding letter: %v", loc, err)
		return errWhenAddingLetter(letter.Body[:len(letter.Body)%10], err)
	}

	s.log.Infof("%s: successfully saved letter", loc)
	return nil
}

func (s *sqliteDB) Delete(id int) error {
	loc := GLOC + "Delete()"
	s.log.Debugf("%s: preparing to delete letter with id: %d", loc, id)

	query := `DELETE FROM letters WHERE id = ?`
	s.log.Debugf("%s: executing query: %s with id=%d", loc, query, id)
	result, err := s.db.Exec(query, id)
	if err != nil {
		s.log.Errorf("%s: error deleting letter with id %d: %v", loc, id, err)
		return errCannotDeleteLetter(id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		s.log.Errorf("%s: error fetching rows affected for id %d: %v", loc, id, err)
		return errCannotFetchRows(id, err)
	}
	s.log.Debugf("%s: rows affected: %d", loc, rowsAffected)

	if rowsAffected == 0 {
		s.log.Warnf("%s: no letter found with id %d", loc, id)
		return errLetterNotFound(id)
	}

	s.log.Infof("%s: successfully deleted letter with id %d", loc, id)
	return nil
}

func (s *sqliteDB) GetAll() ([]st.Letter, error) {
	loc := GLOC + "GetAll()"
	s.log.Debugf("%s: preparing to get all letters", loc)

	query := `SELECT l.id, l.body, l.found_at, l.found_in, 
	          TRIM(a.fname || ' ' || a.mname || ' ' || a.lname) AS author
              FROM letters l
              LEFT JOIN authors a ON l.author_id = a.id`
	s.log.Debugf("%s: executing query: %s", loc, query)
	rows, err := s.db.Query(query)
	if err != nil {
		s.log.Errorf("%s: error executing query: %v", loc, err)
		return nil, err
	}
	defer func() {
		s.log.Debugf("%s: closing rows", loc)
		rows.Close()
	}()

	var letters []st.Letter
	for rows.Next() {
		var letter st.Letter
		err := rows.Scan(
			&letter.ID,
			&letter.Body,
			&letter.FoundAt,
			&letter.FoundIn,
			&letter.Author,
		)
		if err != nil {
			s.log.Errorf("%s: error scanning row: %v", loc, err)
			return nil, err
		}
		s.log.Debugf("%s: scanned letter: %+v", loc, letter)
		letters = append(letters, letter)
	}

	if err = rows.Err(); err != nil {
		s.log.Errorf("%s: error iterating rows: %v", loc, err)
		return nil, err
	}

	s.log.Infof("%s: successfully retrieved %d letters", loc, len(letters))
	return letters, nil
}

// getOrCreateAuthor checks if author exists in db, if not creates it
func (s *sqliteDB) getOrCreateAuthor(name string) (int, error) {
	loc := GLOC + "getOrCreateAuthor()"
	s.log.Debugf("%s: splitting author name: '%s'", loc, name)

	var fname, mname, lname string = lib.SplitName(name)
	s.log.Debugf("%s: split result: fname='%s', mname='%s', lname='%s'", loc, fname, mname, lname)
	var id int

	query := `SELECT id, fname, mname, lname
	          FROM authors
			  WHERE fname = ? AND mname = ? AND lname = ?`
	s.log.Debugf("%s: executing query: %s with fname=%s, mname=%s, lname=%s", loc, query, fname, mname, lname)
	if err := s.db.QueryRow(query, fname, mname, lname).Scan(
		&id,
		&fname,
		&mname,
		&lname,
	); err == nil {
		s.log.Debugf("%s: author found: id=%d, fname=%s, mname=%s, lname=%s", loc, id, fname, mname, lname)
		s.log.Infof("%s: author '%s' found with id %d", loc, name, id)
		return id, nil
	} else if err != sql.ErrNoRows {
		s.log.Errorf("%s: error fetching author '%s': %v", loc, name, err)
		return -1, errWnenFetchAuthor(name, err)
	}

	s.log.Debugf("%s: author not found, inserting new author: fname=%s, mname=%s, lname=%s", loc, fname, mname, lname)
	query = `INSERT INTO authors (fname, mname, lname)
				VALUES (?, ?, ?)`
	res, err := s.db.Exec(query, fname, mname, lname)
	if err != nil {
		s.log.Errorf("%s: error adding author '%s': %v", loc, name, err)
		return -1, errWhenAddingAuthor(name, err)
	}

	uid, err := res.LastInsertId()
	if err != nil {
		s.log.Errorf("%s: error getting last insert id for author '%s': %v", loc, name, err)
		return -1, err
	}

	s.log.Infof("%s: successfully added author '%s' with id %d", loc, name, uid)
	return int(uid), nil
}
