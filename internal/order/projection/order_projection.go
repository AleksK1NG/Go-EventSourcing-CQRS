package projection

import (
	"context"
	"github.com/AleksK1NG/es-microservice/pkg/logger"
	"github.com/EventStore/EventStore-Client-Go/esdb"
)

type orderProjection struct {
	log logger.Logger
	db  *esdb.Client
}

func NewOrderProjection(log logger.Logger, db *esdb.Client) *orderProjection {
	return &orderProjection{log: log, db: db}
}

func (o *orderProjection) ProcessEvents(ctx context.Context) error {
	o.log.Info("starting order subscription")

	stream, err := o.db.SubscribeToAll(ctx, esdb.SubscribeToAllOptions{
		Filter: &esdb.SubscriptionFilter{
			Type:     esdb.StreamFilterType,
			Prefixes: []string{"order-"},
		},
	})
	if err != nil {
		return err
	}
	defer stream.Close()

	for {
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
			o.log.Infof("subscription event StreamID: %s", event.EventAppeared.Event.StreamID)
			o.log.Infof("subscription event Data: %s", string(event.EventAppeared.Event.Data))
		}
	}

	return nil
}
