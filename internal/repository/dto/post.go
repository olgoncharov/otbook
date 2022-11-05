package dto

import "time"

type PostShortInfo struct {
	ID        uint64    `json:"id"`
	Author    string    `json:"author"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
}

type PostFilters struct {
	Authors  []string
	DateFrom *time.Time
}
