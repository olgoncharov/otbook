package profileslist

import (
	"github.com/olgoncharov/otbook/internal/pkg/types"
)

type profileInfo struct {
	Username  string     `json:"username"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Birthdate types.Date `json:"birthdate"`
	City      string     `json:"city"`
	Sex       string     `json:"sex"`
	Hobby     string     `json:"hobby"`
}

type response struct {
	List       []profileInfo `json:"list"`
	TotalCount uint          `json:"totalCount"`
}
