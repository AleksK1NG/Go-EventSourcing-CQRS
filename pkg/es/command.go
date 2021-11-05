package es

// Command commands interface for event sourcing.
type Command interface {
	GetAggregateID() string
}
