package storage

import "time"

type Letter struct {
	ID      int       `json:"id"`
	Author  string    `json:"author"`
	FoundAt time.Time `json:"found_at"`
	FoundIn string    `json:"found_in"`
	Body    string    `json:"body"`
}

type LettersStorage interface {
	Get(id int) (Letter, error)
	Save(letter Letter) error
	Delete(id int) error
	Update(letter Letter) error

	GetAll() ([]Letter, error)
}
