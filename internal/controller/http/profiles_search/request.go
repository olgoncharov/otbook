package profilessearch

import "github.com/olgoncharov/otbook/internal/pkg/errgroup"

type requestBody struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (r *requestBody) validate() error {
	eg := errgroup.NewErrorGroup("\n")

	if r.FirstName == "" {
		eg.AddErrorText("empty firstName")
	}

	if r.LastName == "" {
		eg.AddErrorText("empty lastName")
	}

	return eg.Err()
}
