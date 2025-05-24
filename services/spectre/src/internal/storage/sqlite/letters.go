package sqlite

// ! TODO : think about copy-paste

import (
	"database/sql"
	"spectre/internal/models"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const (
	GLOC_LTS = "src/internal/storage/sqlite/letters.go/" // for logging
)

// GetLetterByID retrieves a letter by its ID from the database.
func (s *sqliteDB) GetLetterByID(id int) (models.Letter, error) {
	loc := GLOC_LTS + "GetLetterByID()"
	s.log.Infof("%s: called for id=%d", loc, id)

	query := `SELECT l.id, l.body, l.found_at, l.found_in
              FROM letters l
              LEFT JOIN authors a ON l.author_id = a.id
              WHERE l.id = ?`
	s.log.Debugf("%s: query: %s", loc, query)

	var letter models.Letter
	if err := s.db.QueryRow(query, id).Scan(
		&letter.ID,
		&letter.Body,
		&letter.FoundAt,
		&letter.FoundIn,
	); err != nil {
		if err == sql.ErrNoRows {
			s.log.Warnf("%s: letter not found with id %d", loc, id)
			return models.Letter{}, errLetterNotFound(id)
		}
		s.log.Errorf("%s: error retrieving letter with id %d: %v", loc, id, err)
		return models.Letter{}, errCannotGetLetter(id, err)
	}

	s.log.Infof("%s: letter retrieved id=%d", loc, id)
	return letter, nil
}

// SaveLetter saves a letter to the database.
func (s *sqliteDB) SaveLetter(letter models.Letter) error {
	loc := GLOC_LTS + "SaveLetter()"
	s.log.Infof("%s: called for author='%s'", loc, letter.Author)

	aID, err := s.getOrCreateAuthor(letter.Author)
	if err != nil || aID <= 0 {
		s.log.Errorf("%s: error getting or creating author '%s': %v", loc, letter.Author, err)
		return err
	}

	query := `INSERT INTO letters (body, found_at, found_in, author_id)
	          VALUES (?, ?, ?, ?)`
	s.log.Debugf("%s: query: %s", loc, query)
	_, err = s.db.Exec(query,
		letter.Body,
		letter.FoundAt,
		letter.FoundIn,
		aID)
	if err != nil {
		s.log.Errorf("%s: error adding letter: %v", loc, err)
		return errWhenAddingLetter(letter.Body[:len(letter.Body)%10], err)
	}

	s.log.Infof("%s: letter saved", loc)
	return nil
}

// DeleteLetter deletes a letter by its ID from the database.
func (s *sqliteDB) DeleteLetter(id int) error {
	loc := GLOC_LTS + "DeleteLetter()"
	s.log.Infof("%s: called for id=%d", loc, id)

	query := `DELETE FROM letters 
			  WHERE id = ?`
	s.log.Debugf("%s: query: %s", loc, query)
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

	if rowsAffected == 0 {
		s.log.Warnf("%s: no letter found with id %d", loc, id)
		return errLetterNotFound(id)
	}

	s.log.Infof("%s: letter deleted id=%d", loc, id)
	return nil
}

// UpdateLetter updates an existing letter in the database.
func (s *sqliteDB) UpdateLetter(letter models.Letter) error {
	loc := GLOC_LTS + "UpdateLetter()"
	s.log.Infof("%s: called for id=%d", loc, letter.ID)

	aID, err := s.getOrCreateAuthor(letter.Author)
	if err != nil || aID <= 0 {
		s.log.Errorf("%s: error getting or creating author '%s': %v", loc, letter.Author, err)
		return err
	}

	query := `UPDATE letters
              SET body = ?, found_at = ?, found_in = ?, author_id = ?
              WHERE id = ?`
	s.log.Debugf("%s: query: %s", loc, query)
	result, err := s.db.Exec(query,
		letter.Body,
		letter.FoundAt,
		letter.FoundIn,
		aID,
		letter.ID,
	)
	if err != nil {
		s.log.Errorf("%s: error updating letter id=%d: %v", loc, letter.ID, err)
		return errWhenUpdateLetter(letter.Body[:len(letter.Body)%10], err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		s.log.Errorf("%s: error fetching rows affected for id %d: %v", loc, letter.ID, err)
		return err
	}
	if rowsAffected == 0 {
		s.log.Warnf("%s: no letter found with id %d", loc, letter.ID)
		return errLetterNotFound(letter.ID)
	}

	s.log.Infof("%s: letter updated id=%d", loc, letter.ID)
	return nil
}

// GetAllLettersWithAccess retrieves all letters with access level less than or equal to the specified value.
func (s *sqliteDB) GetAllLettersWithAccess(accessLevel int) ([]models.Letter, error) {
	loc := GLOC_LTS + "GetAllLettersWithAccess()"
	s.log.Infof("%s: called with accessLevel=%d", loc, accessLevel)

	query := `SELECT l.id, l.body, l.found_at, l.found_in, 
	          TRIM(a.fname || ' ' || a.mname || ' ' || a.lname) AS author
              FROM letters l
              LEFT JOIN authors a ON l.author_id = a.id
			  WHERE l.access_level <= ?`
	s.log.Debugf("%s: query: %s", loc, query)
	rows, err := s.db.Query(query, accessLevel)
	if err != nil {
		s.log.Errorf("%s: error executing query: %v", loc, err)
		return nil, err
	}
	defer rows.Close()

	var letters []models.Letter
	for rows.Next() {
		var letter models.Letter
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
		letters = append(letters, letter)
	}

	if err = rows.Err(); err != nil {
		s.log.Errorf("%s: error iterating rows: %v", loc, err)
		return nil, err
	}

	s.log.Infof("%s: retrieved %d letters", loc, len(letters))
	return letters, nil
}

// GetAllLetters retrieves all letters from the database.
func (s *sqliteDB) GetAllLetters() ([]models.Letter, error) {
	loc := GLOC_LTS + "GetAllLetters)"
	s.log.Infof("%s: called", loc)

	query := `SELECT l.id, l.body, l.found_at, l.found_in, 
	          TRIM(a.fname || ' ' || a.mname || ' ' || a.lname) AS author
              FROM letters l
              LEFT JOIN authors a ON l.author_id = a.id`
	s.log.Debugf("%s: query: %s", loc, query)
	rows, err := s.db.Query(query)
	if err != nil {
		s.log.Errorf("%s: error executing query: %v", loc, err)
		return nil, err
	}
	defer rows.Close()

	var letters []models.Letter
	for rows.Next() {
		var letter models.Letter
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
		letters = append(letters, letter)
	}

	if err = rows.Err(); err != nil {
		s.log.Errorf("%s: error iterating rows: %v", loc, err)
		return nil, err
	}

	s.log.Infof("%s: retrieved %d letters", loc, len(letters))
	return letters, nil
}

// splitName splits a full name string into first, middle, and last name parts.
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

// getOrCreateAuthor checks if an author exists in the database, and creates one if not. Returns the author ID.
func (s *sqliteDB) getOrCreateAuthor(name string) (int, error) {
	loc := GLOC_LTS + "getOrCreateAuthor()"
	s.log.Infof("%s: called for name='%s'", loc, name)

	var fname, mname, lname string = splitName(name)
	var id int

	query := `SELECT id, fname, mname, lname
	          FROM authors
			  WHERE fname = ? AND mname = ? AND lname = ?`
	s.log.Debugf("%s: query: %s", loc, query)
	if err := s.db.QueryRow(query, fname, mname, lname).Scan(
		&id,
		&fname,
		&mname,
		&lname,
	); err == nil {
		s.log.Infof("%s: author found id=%d", loc, id)
		return id, nil
	} else if err != sql.ErrNoRows {
		s.log.Errorf("%s: error fetching author '%s': %v", loc, name, err)
		return -1, errWnenFetchAuthor(name, err)
	}

	query = `INSERT INTO authors (fname, mname, lname)
				VALUES (?, ?, ?)`
	s.log.Debugf("%s: query: %s", loc, query)
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

	s.log.Infof("%s: author added id=%d", loc, uid)
	return int(uid), nil
}
