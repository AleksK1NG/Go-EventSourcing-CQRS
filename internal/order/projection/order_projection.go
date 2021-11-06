package projection

import (
	"context"
	"github.com/AleksK1NG/es-microservice/internal/order/events"
	"github.com/AleksK1NG/es-microservice/internal/order/repository"
	"github.com/AleksK1NG/es-microservice/pkg/es"
	"github.com/AleksK1NG/es-microservice/pkg/logger"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"strings"
)

type orderProjection struct {
	log       logger.Logger
	db        *esdb.Client
	mongoRepo repository.OrderRepository
}

func NewOrderProjection(log logger.Logger, db *esdb.Client, mongoRepo repository.OrderRepository) *orderProjection {
	return &orderProjection{log: log, db: db, mongoRepo: mongoRepo}
}

type Worker func(ctx context.Context, stream *esdb.Subscription, workerID int) error

func (o *orderProjection) Subscribe(ctx context.Context, prefixes []string, poolSize int, worker Worker) error {
	o.log.Infof("starting order subscription: %+v", prefixes)

	stream, err := o.db.SubscribeToAll(ctx, esdb.SubscribeToAllOptions{
		Filter: &esdb.SubscriptionFilter{Type: esdb.StreamFilterType, Prefixes: prefixes},
	})
	if err != nil {
		return err
	}
	defer stream.Close()

	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i <= poolSize; i++ {
		g.Go(func() error {
			return worker(ctx, stream, i)
		})
	}
	return g.Wait()
}

func (o *orderProjection) ProcessEvents(ctx context.Context, stream *esdb.Subscription, workerID int) error {

	for {
		select {
		case <-ctx.Done():
			o.log.Errorf("ctxDone: %v", ctx.Err())
			return ctx.Err()
		default:
		}

		event := stream.Recv()

		if event.SubscriptionDropped != nil {
			o.log.Error("Subscription Dropped")
			if event.SubscriptionDropped.Error != nil {
				o.log.Errorf("SubscriptionDropped error: %s", event.SubscriptionDropped.Error.Error())
				return event.SubscriptionDropped.Error
			}
			return errors.New("Subscription Dropped")
		}

		if event.EventAppeared != nil {
			streamId := event.EventAppeared.OriginalEvent().StreamID
			revision := event.EventAppeared.OriginalEvent().EventNumber

			o.log.Infof("received event %v@%v", revision, streamId)
			//o.log.Infof("subscription event StreamID: %s", event.EventAppeared.Event.StreamID)
			//o.log.Infof("subscription event Data: %s", string(event.EventAppeared.Event.Data))
			o.log.Infof("process subscription: %s workerID: %v", stream.Id(), workerID)

			err := o.When(ctx, es.NewEventFromRecorded(event.EventAppeared.Event))
			if err != nil {
				o.log.Errorf("order projection when: %v", err)
			}
		}
	}
}

func (o *orderProjection) When(ctx context.Context, evt es.Event) error {
	switch evt.GetEventType() {

	case events.OrderCreated:
		return o.handleOrderCreateEvent(ctx, evt)

	case events.OrderPaid:
		return o.handleOrderPaidEvent(ctx, evt)

	case events.OrderSubmitted:
		return o.handleSubmitEvent(ctx, evt)

	case events.OrderDelivering:
		return nil

	case events.OrderDelivered:
		return nil

	case events.OrderCanceled:
		return nil

	case events.OrderUpdated:
		return o.handleUpdateEvent(ctx, evt)
	default:
		return es.ErrInvalidEventType
	}
}

func GetOrderAggregateID(eventAggregateID string) string {
	return strings.ReplaceAll(eventAggregateID, "order-", "")
}
