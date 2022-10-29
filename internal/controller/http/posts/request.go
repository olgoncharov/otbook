package posts

import "github.com/olgoncharov/otbook/internal/pkg/errgroup"

type createRequest struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func (r *createRequest) validate() error {
	eg := errgroup.NewErrorGroup("\n")

	if r.Title == "" {
		eg.AddErrorText("empty title")
	}

	if len(r.Title) > 250 {
		eg.AddErrorText("title length must be less than 250 symbols")
	}

	if r.Text == "" {
		eg.AddErrorText("empty text")
	}

	return eg.Err()
}
