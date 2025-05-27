package models

type User struct {
	ID          int    `json:"id"`
	Login       string `json:"login"`
	AccessLevel int    `json:"access_level"`
	PHash       []byte
}
