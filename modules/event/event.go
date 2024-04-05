package event

import "time"

type Event struct {
	Timestamp time.Time
	Id        uint64
	Version   uint64
}
