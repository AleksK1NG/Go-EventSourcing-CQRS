package store

import (
	"context"
	"github.com/AleksK1NG/es-microservice/pkg/es"
	"github.com/AleksK1NG/es-microservice/pkg/logger"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/pkg/errors"
	"io"
)

type aggregateStore struct {
	log logger.Logger
	db  *esdb.Client
}

func NewAggregateStore(log logger.Logger, db *esdb.Client) *aggregateStore {
	return &aggregateStore{log: log, db: db}
}

func (a *aggregateStore) Load(ctx context.Context, aggregate es.Aggregate) error {
	stream, err := a.db.ReadStream(ctx, aggregate.GetID(), esdb.ReadStreamOptions{
		//Direction: esdb.Forwards,
		//From:      esdb.Revision(1),
	}, 100)
	if err != nil {
		return err
	}
	defer stream.Close()

	events := make([]es.Event, 0, 100)
	for {
		event, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		events = append(events, es.NewEventFromRecorded(event.Event))
	}

	return aggregate.Load(events)
}

func (a *aggregateStore) Save(ctx context.Context, aggregate es.Aggregate) error {
	eventsData := make([]esdb.EventData, 0, len(aggregate.GetUncommittedEvents()))
	for _, event := range aggregate.GetUncommittedEvents() {
		eventsData = append(eventsData, event.ToEventData())
	}

	var expectedRevision esdb.ExpectedRevision = esdb.StreamExists{}
	if aggregate.GetVersion() == 1 {
		expectedRevision = esdb.NoStream{}
	}
	a.log.Infof("SaveEvents expectedRevision: %T", expectedRevision)

	stream, err := a.db.AppendToStream(ctx, aggregate.GetID(), esdb.AppendToStreamOptions{
		ExpectedRevision: expectedRevision,
	}, eventsData...)
	if err != nil {
		return err
	}

	a.log.Infof("SaveEvents stream: %+v", stream)
	return nil
}

func (a *aggregateStore) Exists(ctx context.Context, streamID string) error {
	stream, err := a.db.ReadStream(ctx, streamID, esdb.ReadStreamOptions{
		Direction: esdb.Backwards,
		From:      esdb.Revision(1),
	}, 1)
	if err != nil {
		return err
	}
	defer stream.Close()

	for {
		_, err := stream.Recv()
		if errors.Is(err, esdb.ErrStreamNotFound) {
			return esdb.ErrStreamNotFound
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
	}

	return nil
}
