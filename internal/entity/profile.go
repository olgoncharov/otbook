package entity

import "time"

type Profile struct {
	Username    string
	FirstName   string
	LastName    string
	Birthdate   time.Time
	City        string
	Sex         string
	Hobby       string
	IsCelebrity bool
}
