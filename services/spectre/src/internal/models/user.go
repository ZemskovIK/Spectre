package models

type User struct {
	ID          int    `json:"id"`
	Login       string `json:"login"`
	PHash       []byte
	AccessLevel int `json:"access_level"`
}
