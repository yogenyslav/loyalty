package model

import "time"

// PollableOrder represents an order that can be polled for updates.
type PollableOrder struct {
	Number     string    `db:"number"`
	Attempt    int       `db:"attempt"`
	LastPolled time.Time `db:"last_polled"`
	NextPoll   time.Time `db:"next_poll"`
}
