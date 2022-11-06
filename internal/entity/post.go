package entity

import "time"

type Post struct {
	ID        uint64
	Author    string
	Title     string
	Text      string
	CreatedAt time.Time
}
