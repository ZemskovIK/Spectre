package models

type User struct {
	ID          int    `json:"id"`
	Login       string `json:"login"`
	PHash       []byte `json:"pass_hash"`
	AccessLevel int    `json:"access_level"`
}
