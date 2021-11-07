package mongo_projection

import (
	"context"
	"github.com/AleksK1NG/es-microservice/config"
	"github.com/AleksK1NG/es-microservice/internal/order/events"
	"github.com/AleksK1NG/es-microservice/internal/order/repository"
	"github.com/AleksK1NG/es-microservice/pkg/constants"
	"github.com/AleksK1NG/es-microservice/pkg/es"
	"github.com/AleksK1NG/es-microservice/pkg/logger"
	"github.com/AleksK1NG/es-microservice/pkg/tracing"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type orderProjection struct {
	log       logger.Logger
	db        *esdb.Client
	cfg       *config.Config
	mongoRepo repository.OrderRepository
}

func NewOrderProjection(log logger.Logger, db *esdb.Client, mongoRepo repository.OrderRepository, cfg *config.Config) *orderProjection {
	return &orderProjection{log: log, db: db, mongoRepo: mongoRepo, cfg: cfg}
}

type Worker func(ctx context.Context, stream *esdb.PersistentSubscription, workerID int) error

func (o *orderProjection) Subscribe(ctx context.Context, prefixes []string, poolSize int, worker Worker) error {
	o.log.Infof("starting order subscription: %+v", prefixes)

	err := o.db.CreatePersistentSubscriptionAll(ctx, o.cfg.Subscriptions.MongoProjectionGroupName, esdb.PersistentAllSubscriptionOptions{
		Filter: &esdb.SubscriptionFilter{Type: esdb.StreamFilterType, Prefixes: prefixes},
	})
	if err != nil {
		if subscriptionError, ok := err.(*esdb.PersistentSubscriptionError); !ok || ok && (subscriptionError.Code != 6) {
			o.log.Errorf("CreatePersistentSubscriptionAll: %v", subscriptionError.Error())
		}
	}

	stream, err := o.db.ConnectToPersistentSubscription(ctx, constants.EsAll, o.cfg.Subscriptions.MongoProjectionGroupName, esdb.ConnectToPersistentSubscriptionOptions{})
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

func (o *orderProjection) ProcessEvents(ctx context.Context, stream *esdb.PersistentSubscription, workerID int) error {

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
			streamID := event.EventAppeared.OriginalEvent().StreamID
			revision := event.EventAppeared.OriginalEvent().EventNumber
			o.log.Infof("(mongo projection event): revision: %v, streamID: %v, workerID: %v, eventType: %s", revision, streamID, workerID, event.EventAppeared.Event.EventType)

			err := o.When(ctx, es.NewEventFromRecorded(event.EventAppeared.Event))
			if err != nil {
				o.log.Errorf("order projection when: %v", err)
				if err := stream.Nack(err.Error(), esdb.Nack_Unknown, event.EventAppeared); err != nil {
					o.log.Errorf("stream.Nack: %v", err)
					return err
				}
			}
			err = stream.Ack(event.EventAppeared)
			if err != nil {
				o.log.Errorf("stream.Ack: %v", err)
				return err
			}
			o.log.Infof("(ACK event commit): %v", *event.EventAppeared.Commit)
		}
	}
}

func (o *orderProjection) When(ctx context.Context, evt es.Event) error {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "orderProjection.When")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()), log.String("EventType", evt.GetEventType()))

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
		o.log.Debugf("when eventType: %s", evt.EventType)
		return nil
	}
}
