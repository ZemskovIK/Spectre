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

	GetUserByLogin(login string) (models.User, error)
	SaveUser(id int) error
	DeleteUser(id int) error
	UpdateUser(usr models.User) error
}

// asd.com:PORT
// PORT /enctypt_bytes POST -> json
```
{
  "content": { ... },   // полезная нагрузка (зашифрованные данные)
  "error": null | string, // описание ошибки, если она есть
  "iv": "base64",       // IV (инициализационный вектор) для расшифровки
  "hmac": "base64",     // HMAC для проверки целостности
  "nonce": "base64"     // nonce, если используется
}
```