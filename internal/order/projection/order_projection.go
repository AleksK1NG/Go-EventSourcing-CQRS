package projection

import (
	"context"
	"github.com/AleksK1NG/es-microservice/internal/models"
	"github.com/AleksK1NG/es-microservice/internal/order/events"
	"github.com/AleksK1NG/es-microservice/internal/order/repository"
	"github.com/AleksK1NG/es-microservice/pkg/es"
	"github.com/AleksK1NG/es-microservice/pkg/logger"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"strings"
	"sync"
)

type orderProjection struct {
	log       logger.Logger
	db        *esdb.Client
	mongoRepo repository.OrderRepository
}

func NewOrderProjection(log logger.Logger, db *esdb.Client, mongoRepo repository.OrderRepository) *orderProjection {
	return &orderProjection{log: log, db: db, mongoRepo: mongoRepo}
}

type Worker func(ctx context.Context, stream *esdb.Subscription, wg *sync.WaitGroup, workerID int)

func (o *orderProjection) Subscribe(ctx context.Context, prefixes []string, poolSize int, worker Worker) error {
	o.log.Infof("starting order subscription prefixes: %+v", prefixes)

	stream, err := o.db.SubscribeToAll(ctx, esdb.SubscribeToAllOptions{
		Filter: &esdb.SubscriptionFilter{Type: esdb.StreamFilterType, Prefixes: prefixes},
	})
	if err != nil {
		return err
	}
	defer stream.Close()

	wg := &sync.WaitGroup{}
	for i := 0; i <= poolSize; i++ {
		wg.Add(1)
		go worker(ctx, stream, wg, i)
	}

	wg.Wait()
	return nil
}

func (o *orderProjection) ProcessEvents(ctx context.Context, stream *esdb.Subscription, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			o.log.Errorf("ctxDone: %v", ctx.Err())
			return
		default:
		}

		event := stream.Recv()

		if event.SubscriptionDropped != nil {
			o.log.Error("Subscription Dropped")
			if event.SubscriptionDropped.Error != nil {
				o.log.Errorf("SubscriptionDropped error: %s", event.SubscriptionDropped.Error.Error())
			}
			return
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

	o.log.Infof("subscription finished: %sm workerID: %v", stream.Id(), workerID)
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

func (o *orderProjection) handleOrderCreateEvent(ctx context.Context, evt es.Event) error {
	var eventData events.OrderCreatedData
	if err := evt.GetJsonData(&eventData); err != nil {
		return err
	}

	op := &models.OrderProjection{
		OrderID:    GetOrderAggregateID(evt.AggregateID),
		ItemsIDs:   eventData.ItemsIDs,
		Created:    true,
		Paid:       false,
		Submitted:  false,
		Delivering: false,
		Delivered:  false,
		Canceled:   false,
	}

	result, err := o.mongoRepo.Insert(ctx, op)
	if err != nil {
		return err
	}

	o.log.Debugf("projection OrderCreated result: %s", result)
	return nil
}

func (o *orderProjection) handleOrderPaidEvent(ctx context.Context, evt es.Event) error {
	op := &models.OrderProjection{OrderID: GetOrderAggregateID(evt.AggregateID), Paid: true}
	return o.mongoRepo.UpdateOrder(ctx, op)
}

func (o *orderProjection) handleSubmitEvent(ctx context.Context, evt es.Event) error {
	op := &models.OrderProjection{OrderID: GetOrderAggregateID(evt.AggregateID), Submitted: true}
	return o.mongoRepo.UpdateOrder(ctx, op)
}

func (o *orderProjection) handleUpdateEvent(ctx context.Context, evt es.Event) error {
	var eventData events.OrderCreatedData
	if err := evt.GetJsonData(&eventData); err != nil {
		return err
	}

	op := &models.OrderProjection{OrderID: GetOrderAggregateID(evt.AggregateID), ItemsIDs: eventData.ItemsIDs}
	return o.mongoRepo.UpdateOrder(ctx, op)
}

func GetOrderAggregateID(eventAggregateID string) string {
	return strings.ReplaceAll(eventAggregateID, "order-", "")
}
