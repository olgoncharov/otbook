package dto

import "time"

type PostShortInfo struct {
	ID        uint64
	Author    string
	Title     string
	CreatedAt time.Time
}
