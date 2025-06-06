package storage

import (
	"spectre/internal/models"
)

type Storage interface {
	GetLetterByID(id int) (models.Letter, error)
	SaveLetter(letter models.Letter) error
	DeleteLetter(id int) error
	UpdateLetter(letter models.Letter) error
	GetAllLettersWithAccess(accessLevel int) ([]models.Letter, error)
	GetAllLetters() ([]models.Letter, error)

	GetUserByLogin(login string) (models.User, error)
	SaveUser(usr models.User) error
	DeleteUser(id int) error
	UpdateUser(usr models.User) error
	GetAllUsers() ([]models.User, error)
	GetUserByID(id int) (models.User, error)
}
