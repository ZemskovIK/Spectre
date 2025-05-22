package storage

import (
	"spectre/internal/models"
)

type Storage interface {
	Get(id int) (models.Letter, error)
	Save(letter models.Letter) error
	Delete(id int) error
	Update(letter models.Letter) error
	GetAll() ([]models.Letter, error)

	GetUserByLogin(login string) (models.User, error)
}
