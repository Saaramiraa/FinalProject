package entities

import "time"

type Author struct {
	Id   int
	Name string
	DoB  string
	Updated_At time.Time
}