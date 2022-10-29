package posts

import "time"

type (
	postInfo struct {
		ID        uint64            `json:"id"`
		Author    string            `json:"author"`
		Title     string            `json:"title"`
		CreatedAt time.Time         `json:"createdAt"`
		Links     map[string]string `json:"links"`
	}

	listResponse struct {
		List       []postInfo `json:"list"`
		TotalCount uint       `json:"totalCount"`
	}

	createResponse struct {
		ID uint64 `json:"id"`
	}

	singlePostResponse struct {
		ID        uint64            `json:id"`
		Author    string            `json:"author"`
		Title     string            `json:"title"`
		Text      string            `json:"text"`
		CreatedAt time.Time         `json:"createdAt"`
		Links     map[string]string `json:"links"`
	}
)
