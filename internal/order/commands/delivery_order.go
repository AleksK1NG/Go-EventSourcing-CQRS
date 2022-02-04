package commands

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
	Handle(ctx context.Context, command *aggregate.OrderDeliveredCommand) error
}

type deliveryOrderCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  es.AggregateStore
}

func NewDeliveryOrderCommandHandler(log logger.Logger, cfg *config.Config, es es.AggregateStore) *deliveryOrderCommandHandler {
	return &deliveryOrderCommandHandler{log: log, cfg: cfg, es: es}
}

func (c *deliveryOrderCommandHandler) Handle(ctx context.Context, command *aggregate.OrderDeliveredCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "deliveryOrderCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", command.GetAggregateID()))

	return aggregate.HandleCommand(ctx, c.es, command)
}
