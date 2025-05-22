package models

import "time"

type Letter struct {
	ID      int       `json:"id"`
	Author  string    `json:"author"`
	FoundAt time.Time `json:"found_at"`
	FoundIn string    `json:"found_in"`
	Body    string    `json:"body"`
}
