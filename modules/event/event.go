package event

import "time"

type Event struct {
	Timestamp time.Time
	Type      string
}

type WithEvent struct {
	Event Event
}
