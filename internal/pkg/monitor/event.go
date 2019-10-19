package monitor

import "fmt"

// EventType represents a file event type enumeration
type EventType int

func (e EventType) String() string {
	switch e {
	case Create:
		return "create"
	case Update:
		return "update"
	case Delete:
		return "delete"
	case Check:
		return "check"
	default:
		return "(event)"
	}
}

// EventType enumerations
const (
	Create EventType = iota + 1
	Update
	Delete
	Check
)

// Event represents a file event
type Event struct {
	Filename string
	Type     EventType
}

func (e Event) String() string {
	return fmt.Sprintf("%v: %s", e.Type, e.Filename)
}

// EventChannel represents a channel for Events
type EventChannel chan Event
