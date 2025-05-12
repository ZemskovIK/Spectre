package storage

type Letter struct {
	ID      int    `json:"id"`
	Author  string `json:"author"`
	FoundAt string `json:"found_at"`
	FoundIn string `json:"found_in"`
	Title   string `json:"title"`
	Body    string `json:"body"`
}

type LettersStorage interface {
	Get(id int) (Letter, error)
	Save(letter Letter) error
	Delete(id int) error

	GetAll() ([]Letter, error)
}
