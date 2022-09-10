package myprofile

import (
	"github.com/olgoncharov/otbook/internal/pkg/errgroup"
	"github.com/olgoncharov/otbook/internal/pkg/types"
)

type requestBody struct {
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Birthdate types.Date `json:"birthDate"`
	City      string     `json:"city"`
	Sex       string     `json:"sex"`
	Hobby     string     `json:"hobby"`
}

func (r *requestBody) validate() error {
	eg := errgroup.NewErrorGroup("\n")

	if r.FirstName == "" {
		eg.AddErrorText("empty first name")
	}

	if r.LastName == "" {
		eg.AddErrorText("empty last name")
	}

	if r.Birthdate.IsZero() {
		eg.AddErrorText("empty birthdate")
	}

	if r.City == "" {
		eg.AddErrorText("empty city")
	}

	if r.Sex == "" {
		eg.AddErrorText("empty sex")
	}

	if r.Hobby == "" {
		eg.AddErrorText("empty hobby")
	}

	return eg.Err()
}
