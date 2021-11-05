package store

import (
	"context"
	"github.com/AleksK1NG/es-microservice/pkg/es"
	"github.com/AleksK1NG/es-microservice/pkg/logger"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/pkg/errors"
	"io"
)

type eventStore struct {
	log logger.Logger
	db  *esdb.Client
}

func NewEventStore(log logger.Logger, db *esdb.Client) *eventStore {
	return &eventStore{log: log, db: db}
}

func (e *eventStore) SaveEvents(ctx context.Context, streamID string, events []es.Event) error {
	eventsData := make([]esdb.EventData, 0, len(events))
	for _, event := range events {
		eventsData = append(eventsData, event.ToEventData())
	}

	stream, err := e.db.AppendToStream(ctx, streamID, esdb.AppendToStreamOptions{}, eventsData...)
	if err != nil {
		return err
	}

	e.log.Infof("SaveEvents stream: %+v", stream)
	return nil
}

func (e *eventStore) LoadEvents(ctx context.Context, streamID string) ([]es.Event, error) {
	stream, err := e.db.ReadStream(ctx, streamID, esdb.ReadStreamOptions{
		Direction: esdb.Forwards,
		From:      esdb.Revision(1),
	}, 100)
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	events := make([]es.Event, 0, 100)
	for {
		event, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		events = append(events, es.NewEventFromRecorded(event.Event))
	}

	return events, nil
}
