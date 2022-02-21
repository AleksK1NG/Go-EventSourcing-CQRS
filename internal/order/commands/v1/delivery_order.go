package v1

import (
	"context"

	"github.com/AleksK1NG/es-microservice/config"
	"github.com/AleksK1NG/es-microservice/internal/order/aggregate"
	"github.com/AleksK1NG/es-microservice/pkg/es"
	"github.com/AleksK1NG/es-microservice/pkg/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type DeliveryOrderCommandHandler interface {
	Handle(ctx context.Context, command *OrderDeliveredCommand) error
}

type deliveryOrderCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  es.AggregateStore
}

func NewDeliveryOrderCommandHandler(log logger.Logger, cfg *config.Config, es es.AggregateStore) *deliveryOrderCommandHandler {
	return &deliveryOrderCommandHandler{log: log, cfg: cfg, es: es}
}

func (c *deliveryOrderCommandHandler) Handle(ctx context.Context, command *OrderDeliveredCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "deliveryOrderCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", command.GetAggregateID()))

	order, err := aggregate.LoadOrderAggregate(ctx, c.es, command.GetAggregateID())
	if err != nil {
		return err
	}

	if err := order.DeliverOrder(ctx, command.DeliveryTimestamp); err != nil {
		return err
	}

	return c.es.Save(ctx, order)
}
