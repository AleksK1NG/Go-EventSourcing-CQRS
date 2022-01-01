package aggregate

import (
	"context"
	"github.com/AleksK1NG/es-microservice/internal/models"
	"github.com/AleksK1NG/es-microservice/pkg/es"
	"github.com/AleksK1NG/es-microservice/pkg/tracing"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"strings"
)

func GetShopItemsTotalPrice(shopItems []*models.ShopItem) float64 {
	var totalPrice float64 = 0
	for _, item := range shopItems {
		totalPrice += item.Price * float64(item.Quantity)
	}
	return totalPrice
}

// GetOrderAggregateID get order aggregate id for eventstoredb
func GetOrderAggregateID(eventAggregateID string) string {
	return strings.ReplaceAll(eventAggregateID, "order-", "")
}

func IsAggregateNotFound(aggregate es.Aggregate) bool {
	return len(aggregate.GetAppliedEvents()) == 0 || aggregate.GetVersion() == 0
}

// HandleCommand check exists, Load es.Aggregate, HandleCommand and Save to event store
func HandleCommand(ctx context.Context, eventStore es.AggregateStore, command es.Command) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HandleCommandWithExists")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", command.GetAggregateID()))

	order := NewOrderAggregateWithID(command.GetAggregateID())
	if err := eventStore.Load(ctx, order); err != nil {
		return err
	}

	if err := order.HandleCommand(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	span.LogFields(log.String("order", order.Order.String()))
	return eventStore.Save(ctx, order)
}

// HandleCommandWithExists check exists, Load es.Aggregate, HandleCommand and Save to event store
func HandleCommandWithExists(ctx context.Context, eventStore es.AggregateStore, command es.Command) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HandleCommandWithExists")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", command.GetAggregateID()))

	err := eventStore.Exists(ctx, command.GetAggregateID())
	if err != nil && !errors.Is(err, esdb.ErrStreamNotFound) {
		return err
	}

	order := NewOrderAggregateWithID(command.GetAggregateID())
	if err := eventStore.Load(ctx, order); err != nil {
		return err
	}

	if err := order.HandleCommand(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	span.LogFields(log.String("order", order.Order.String()))
	return eventStore.Save(ctx, order)
}