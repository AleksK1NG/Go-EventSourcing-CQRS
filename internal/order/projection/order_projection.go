package projection

import (
	"context"
	"github.com/AleksK1NG/es-microservice/pkg/logger"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"sync"
)

type orderProjection struct {
	log logger.Logger
	db  *esdb.Client
}

func NewOrderProjection(log logger.Logger, db *esdb.Client) *orderProjection {
	return &orderProjection{log: log, db: db}
}

// Worker kafka consumer worker fetch and process messages from reader
type Worker func(ctx context.Context, stream *esdb.Subscription, wg *sync.WaitGroup, workerID int)

func (o *orderProjection) Subscribe(ctx context.Context, prefixes []string, poolSize int, worker Worker) error {
	o.log.Info("starting order subscription")

	//Prefixes: []string{"order-"},
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
	//o.log.Info("starting order subscription")
	//
	//stream, err := o.db.SubscribeToAll(ctx, esdb.SubscribeToAllOptions{
	//	Filter: &esdb.SubscriptionFilter{
	//		Type:     esdb.StreamFilterType,
	//		Prefixes: []string{"order-"},
	//	},
	//})
	//if err != nil {
	//	return err
	//}
	//defer stream.Close()
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
			stream.Close()
			break
		}

		if event.EventAppeared != nil {
			streamId := event.EventAppeared.OriginalEvent().StreamID
			revision := event.EventAppeared.OriginalEvent().EventNumber

			o.log.Infof("received event %v@%v", revision, streamId)
			//o.log.Infof("subscription event StreamID: %s", event.EventAppeared.Event.StreamID)
			//o.log.Infof("subscription event Data: %s", string(event.EventAppeared.Event.Data))
			o.log.Infof("process subscription: %s workerID: %v", stream.Id(), workerID)
		}
	}

	o.log.Infof("subscription finished: %sm workerID: %v", stream.Id(), workerID)
}
