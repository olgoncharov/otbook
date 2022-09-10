package dto

import "time"

type RefreshToken struct {
	Value     string
	CreatedAt time.Time
	TTL       uint64
}
