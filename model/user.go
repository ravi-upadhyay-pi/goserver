package model

import "time"

type User struct {
	Username string
	Password string
	FirstName string
	LastName string
	EmailId string
	PhoneNumber string
	CreatedOn time.Time
	UpdatedOn time.Time
}
