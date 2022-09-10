package signup

import (
	"github.com/olgoncharov/otbook/internal/pkg/errgroup"
	"github.com/olgoncharov/otbook/internal/pkg/types"
)

type requestBody struct {
	Username  string     `json:"username"`
	Password  string     `json:"password"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Birthdate types.Date `json:"birthDate"`
	City      string     `json:"city"`
	Sex       string     `json:"sex"`
	Hobby     string     `json:"hobby"`
}

func (r *requestBody) validate() error {
	eg := errgroup.NewErrorGroup("\n")

	if r.Username == "" {
		eg.AddErrorText("empty username")
	}

	if r.Password == "" {
		eg.AddErrorText("empty password")
	}

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
