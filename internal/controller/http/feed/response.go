package feed

import "time"

type postInfo struct {
	ID        uint64            `json:"id"`
	Author    string            `json:"author"`
	Title     string            `json:"title"`
	CreatedAt time.Time         `json:"createdAt"`
	Links     map[string]string `json:"links"`
}
