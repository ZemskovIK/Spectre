package sqlite

import (
	"database/sql"
	"spectre/internal/models"
)

const (
	GLOC_USR = "src/internal/storage/sqlite/users.go/"
)

func (s *sqliteDB) GetUserByLogin(login string) (models.User, error) {
	loc := GLOC_SQL + "GetUserByLogin()"
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

	return models.User{}, nil
}

func (s *sqliteDB) SaveUser(usr models.User) error {
	loc := GLOC_SQL + "SaveUser()"
	s.log.Debugf("%s: preparing to save user with login: %s", loc, usr.Login)

	query := `INSERT INTO users(login, phash, access_level)
			  VALUES (?, ?, ?)`
	if _, err := s.db.Exec(query,
		usr.Login,
		usr.PHash,
		usr.AccessLevel,
	); err != nil {
		return errCannotSaveUser(usr.Login, err)
	}

	return nil
}

func (s *sqliteDB) DeleteUser(id int) error {
	loc := GLOC_USR + "DeleteUser()"
	_ = loc

	query := `DELETE FROM users
			  WHERE id = ?`
	result, err := s.db.Exec(query, id)
	if err != nil {
		return errCannotDeleteUser(id, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errCannotFetchRows(id, err)
	}
	if rowsAffected == 0 {
		return errUserNotFoundByID(id)
	}

	return nil
}

func (s *sqliteDB) UpdateUser(usr models.User) error {
	query := `UPDATE users 
		      SET login = ?, phash = ?, access_level = ? 
			  WHERE id = ?`
	result, err := s.db.Exec(query,
		usr.Login,
		usr.PHash,
		usr.AccessLevel,
		usr.ID,
	)
	if err != nil {
		return errCannotUpdateUser(usr.ID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errUserNotFoundByID(usr.ID)
	}

	return nil
}

func (s *sqliteDB) GetAllUsers() ([]models.User, error) {
	query := `SELECT * FROM users`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
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
			return nil, err
		}
		users = append(users, usr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
