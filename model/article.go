package model

import "time"

type Article struct {
	Id        string
	Username  string
	Title     string
	Body      string
	Tags      string
	Format    string
	Next      uint64
	Previous  uint64
	Private   bool
	CreatedOn time.Time
	UpdatedOn time.Time
}
