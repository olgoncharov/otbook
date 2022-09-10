package profile

import "github.com/olgoncharov/otbook/internal/pkg/types"

type response struct {
	Username  string     `json:"username"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Birthdate types.Date `json:"birthdate"`
	City      string     `json:"city"`
	Sex       string     `json:"sex"`
	Hobby     string     `json:"hobby"`
}
