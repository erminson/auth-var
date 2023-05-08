package entity

import "time"

type AuthPhone struct {
	Id          int
	Phone       string
	Code        string
	CreatedAt   time.Time
	ValidTil    time.Time
	ConfirmedAt *time.Time
}
