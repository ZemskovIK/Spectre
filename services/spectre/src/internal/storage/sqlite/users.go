package sqlite

import (
	"database/sql"
	"spectre/internal/models"
)

const (
	GLOC_USR = "src/internal/storage/sqlite/users.go/"
)

func (s *sqliteDB) GetUserByLogin(login string) (models.User, error) {
	loc := GLOC_USR + "GetUserByLogin()"
	s.log.Debugf("%s: preparing to get user by login: %s", loc, login)

	query := `SELECT id, login, phash, access_level
			  FROM users
			  WHERE login = ?`
	s.log.Debugf("%s: executing query: %s with login=%s", loc, query, login)

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

	s.log.Infof("%s: successfully retrieved user with login %s", loc, login)
	return user, nil
}

func (s *sqliteDB) SaveUser(usr models.User) error {
	loc := GLOC_USR + "SaveUser()"
	s.log.Debugf("%s: preparing to save user with login: %s", loc, usr.Login)

	query := `INSERT INTO users(login, phash, access_level)
			  VALUES (?, ?, ?)`
	s.log.Debugf("%s: executing query: %s with values: login=%s, phash=****, access_level=%d", loc, query, usr.Login, usr.AccessLevel)
	if _, err := s.db.Exec(query,
		usr.Login,
		usr.PHash,
		usr.AccessLevel,
	); err != nil {
		s.log.Errorf("%s: error saving user with login %s: %v", loc, usr.Login, err)
		return errCannotSaveUser(usr.Login, err)
	}

	s.log.Infof("%s: successfully saved user with login %s", loc, usr.Login)
	return nil
}

func (s *sqliteDB) DeleteUser(id int) error {
	loc := GLOC_USR + "DeleteUser()"
	s.log.Debugf("%s: preparing to delete user with id: %d", loc, id)

	query := `DELETE FROM users
			  WHERE id = ?`
	s.log.Debugf("%s: executing query: %s with id=%d", loc, query, id)
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
	s.log.Debugf("%s: rows affected: %d", loc, rowsAffected)
	if rowsAffected == 0 {
		s.log.Warnf("%s: no user found with id %d", loc, id)
		return errUserNotFoundByID(id)
	}

	s.log.Infof("%s: successfully deleted user with id %d", loc, id)
	return nil
}

func (s *sqliteDB) UpdateUser(usr models.User) error {
	loc := GLOC_USR + "UpdateUser()"
	s.log.Debugf("%s: preparing to update user: %+v", loc, usr)

	query := `UPDATE users 
		      SET login = ?, phash = ?, access_level = ? 
			  WHERE id = ?`
	s.log.Debugf("%s: executing query: %s with values: login=%s, phash=****, access_level=%d, id=%d", loc, query, usr.Login, usr.AccessLevel, usr.ID)
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
	s.log.Debugf("%s: rows affected: %d", loc, rowsAffected)
	if rowsAffected == 0 {
		s.log.Warnf("%s: no user found with id %d", loc, usr.ID)
		return errUserNotFoundByID(usr.ID)
	}

	s.log.Infof("%s: successfully updated user with id %d", loc, usr.ID)
	return nil
}

func (s *sqliteDB) GetAllUsers() ([]models.User, error) {
	loc := GLOC_USR + "GetAllUsers()"
	s.log.Debugf("%s: preparing to get all users", loc)

	query := `SELECT * FROM users`
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
		s.log.Debugf("%s: scanned user: %+v", loc, usr)
		users = append(users, usr)
	}

	if err = rows.Err(); err != nil {
		s.log.Errorf("%s: error iterating rows: %v", loc, err)
		return nil, err
	}

	s.log.Infof("%s: successfully retrieved %d users", loc, len(users))
	return users, nil
}
