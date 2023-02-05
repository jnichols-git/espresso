package events

import "github.com/jakenichols2719/simpleblog/pkg/listings"

type EventKind string

const (
	CreateEvent = EventKind("create")
	UpdateEvent = EventKind("update")
)

type Event struct {
	Kind        EventKind         `json:"event_kind"`
	Timestamp   int64             `json:"event_timestamp"`
	Target      string            `json:"event_target,omitempty"`
	WithContent string            `json:"content,omitempty"`
	WithListing *listings.Listing `json:"listing,omitempty"`
}
