package store

import (
	"context"
	"github.com/AleksK1NG/es-microservice/pkg/es"
	"github.com/AleksK1NG/es-microservice/pkg/logger"
	"github.com/AleksK1NG/es-microservice/pkg/tracing"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"io"
)

const (
	count = 500000
)

type aggregateStore struct {
	log logger.Logger
	db  *esdb.Client
}

func NewAggregateStore(log logger.Logger, db *esdb.Client) *aggregateStore {
	return &aggregateStore{log: log, db: db}
}

func (a *aggregateStore) Load(ctx context.Context, aggregate es.Aggregate) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "aggregateStore.Load")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", aggregate.GetID()))

	stream, err := a.db.ReadStream(ctx, aggregate.GetID(), esdb.ReadStreamOptions{
		//Direction: esdb.Forwards,
		//From:      esdb.Revision(1),
	}, count)
	if err != nil {
		tracing.TraceErr(span, err)
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
			tracing.TraceErr(span, err)
			return err
		}
		events = append(events, es.NewEventFromRecorded(event.Event))
	}

	return aggregate.Load(events)
}

func (a *aggregateStore) Save(ctx context.Context, aggregate es.Aggregate) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "aggregateStore.Save")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", aggregate.GetID()))

	eventsData := make([]esdb.EventData, 0, len(aggregate.GetUncommittedEvents()))
	for _, event := range aggregate.GetUncommittedEvents() {
		eventsData = append(eventsData, event.ToEventData())
	}

	var expectedRevision esdb.ExpectedRevision
	if aggregate.GetVersion() <= 1 {
		expectedRevision = esdb.NoStream{}
		a.log.Infof("SaveEvents expectedRevision: %T", expectedRevision)

		appendStream, err := a.db.AppendToStream(ctx, aggregate.GetID(), esdb.AppendToStreamOptions{ExpectedRevision: expectedRevision}, eventsData...)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		a.log.Infof("SaveEvents stream: %+v", appendStream)
		return nil
	}

	ropts := esdb.ReadStreamOptions{Direction: esdb.Backwards, From: esdb.End{}}
	stream, err := a.db.ReadStream(context.Background(), aggregate.GetID(), ropts, 1)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	defer stream.Close()

	lastEvent, err := stream.Recv()
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	expectedRevision = esdb.Revision(lastEvent.OriginalEvent().EventNumber)
	a.log.Infof("Save expectedRevision: %T", expectedRevision)

	appendStream, err := a.db.AppendToStream(ctx, aggregate.GetID(), esdb.AppendToStreamOptions{ExpectedRevision: expectedRevision}, eventsData...)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	a.log.Infof("Save stream: %+v", appendStream)
	return nil
}

func (a *aggregateStore) Exists(ctx context.Context, streamID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "aggregateStore.Exists")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", streamID))

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
			tracing.TraceErr(span, err)
			return esdb.ErrStreamNotFound
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}
