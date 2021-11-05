package es

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/esdb"
)

// Projection When method works and process Event's like Aggregate's for interacting with read database.
type Projection interface {
	When(ctx context.Context, evt esdb.EventData) error
}
