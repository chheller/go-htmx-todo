package event

import "time"

type Event struct {
	Timestamp time.Time
	Type      string
}
