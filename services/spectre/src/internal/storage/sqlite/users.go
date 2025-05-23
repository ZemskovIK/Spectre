package sqlite

import (
	"database/sql"
	"spectre/internal/models"
)

const (
	GLOC_USR = "src/internal/storage/sqlite/users.go/"
)

// GetUserByLogin retrieves a user by login from the database.
func (s *sqliteDB) GetUserByLogin(login string) (models.User, error) {
	loc := GLOC_USR + "GetUserByLogin()"
	s.log.Infof("%s: called for login=%s", loc, login)

	query := `SELECT id, login, phash, access_level
			  FROM users
			  WHERE login = ?`
	s.log.Debugf("%s: query: %s", loc, query)

	var user models.User
	if err := s.db.QueryRow(query, login).Scan(
		&user.ID,
		&user.Login,
		&user.PHash,
		&user.AccessLevel,
	); err != nil {
		if err == sql.ErrNoRows {
			s.log.Warnf("%s: user not found with login %s", loc, login)
			return models.User{}, errUserNotFoundByLogin(login)
		}
		s.log.Errorf("%s: error retrieving user with login %s: %v", loc, login, err)
		return models.User{}, errCannotGetUser(login, err)
	}

	s.log.Infof("%s: user retrieved login=%s", loc, login)
	return user, nil
}

// SaveUser saves a user to the database.
func (s *sqliteDB) SaveUser(usr models.User) error {
	loc := GLOC_USR + "SaveUser()"
	s.log.Infof("%s: called for login=%s", loc, usr.Login)

	query := `INSERT INTO users(login, phash, access_level)
			  VALUES (?, ?, ?)`
	s.log.Debugf("%s: query: %s", loc, query)
	if _, err := s.db.Exec(query,
		usr.Login,
		usr.PHash,
		usr.AccessLevel,
	); err != nil {
		s.log.Errorf("%s: error saving user with login %s: %v", loc, usr.Login, err)
		return errCannotSaveUser(usr.Login, err)
	}

	s.log.Infof("%s: user saved login=%s", loc, usr.Login)
	return nil
}

// DeleteUser deletes a user by ID from the database.
func (s *sqliteDB) DeleteUser(id int) error {
	loc := GLOC_USR + "DeleteUser()"
	s.log.Infof("%s: called for id=%d", loc, id)

	query := `DELETE FROM users
			  WHERE id = ?`
	s.log.Debugf("%s: query: %s", loc, query)
	result, err := s.db.Exec(query, id)
	if err != nil {
		s.log.Errorf("%s: error deleting user with id %d: %v", loc, id, err)
		return errCannotDeleteUser(id, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		s.log.Errorf("%s: error fetching rows affected for id %d: %v", loc, id, err)
		return errCannotFetchRows(id, err)
	}
	if rowsAffected == 0 {
		s.log.Warnf("%s: no user found with id %d", loc, id)
		return errUserNotFoundByID(id)
	}

	s.log.Infof("%s: user deleted id=%d", loc, id)
	return nil
}

// UpdateUser updates an existing user in the database.
func (s *sqliteDB) UpdateUser(usr models.User) error {
	loc := GLOC_USR + "UpdateUser()"
	s.log.Infof("%s: called for id=%d", loc, usr.ID)

	query := `UPDATE users 
		      SET login = ?, phash = ?, access_level = ? 
			  WHERE id = ?`
	s.log.Debugf("%s: query: %s", loc, query)
	result, err := s.db.Exec(query,
		usr.Login,
		usr.PHash,
		usr.AccessLevel,
		usr.ID,
	)
	if err != nil {
		s.log.Errorf("%s: error updating user with id %d: %v", loc, usr.ID, err)
		return errCannotUpdateUser(usr.ID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		s.log.Errorf("%s: error fetching rows affected for id %d: %v", loc, usr.ID, err)
		return err
	}
	if rowsAffected == 0 {
		s.log.Warnf("%s: no user found with id %d", loc, usr.ID)
		return errUserNotFoundByID(usr.ID)
	}

	s.log.Infof("%s: user updated id=%d", loc, usr.ID)
	return nil
}

// GetAllUsers retrieves all users from the database.
func (s *sqliteDB) GetAllUsers() ([]models.User, error) {
	loc := GLOC_USR + "GetAllUsers()"
	s.log.Infof("%s: called", loc)

	query := `SELECT * FROM users`
	s.log.Debugf("%s: query: %s", loc, query)
	rows, err := s.db.Query(query)
	if err != nil {
		s.log.Errorf("%s: error executing query: %v", loc, err)
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var usr models.User
		err := rows.Scan(
			&usr.ID,
			&usr.Login,
			&usr.PHash,
			&usr.AccessLevel,
		)
		if err != nil {
			s.log.Errorf("%s: error scanning row: %v", loc, err)
			return nil, err
		}
		users = append(users, usr)
	}

	if err = rows.Err(); err != nil {
		s.log.Errorf("%s: error iterating rows: %v", loc, err)
		return nil, err
	}

	s.log.Infof("%s: retrieved %d users", loc, len(users))
	return users, nil
}
