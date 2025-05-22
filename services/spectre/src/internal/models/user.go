package models

type User struct {
	ID          int
	Login       string
	PHash       []byte
	AccessLevel int
}
