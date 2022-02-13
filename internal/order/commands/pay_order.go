package commands

import (
	"context"
	"github.com/AleksK1NG/es-microservice/config"
	"github.com/AleksK1NG/es-microservice/internal/order/aggregate"
	"github.com/AleksK1NG/es-microservice/internal/order/commands/v1"
	"github.com/AleksK1NG/es-microservice/pkg/es"
	"github.com/AleksK1NG/es-microservice/pkg/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type OrderPaidCommandHandler interface {
	Handle(ctx context.Context, command *v1.OrderPaidCommand) error
}

type orderPaidHandler struct {
	log logger.Logger
	cfg *config.Config
	es  es.AggregateStore
}

func NewOrderPaidHandler(log logger.Logger, cfg *config.Config, es es.AggregateStore) *orderPaidHandler {
	return &orderPaidHandler{log: log, cfg: cfg, es: es}
}

func (c *orderPaidHandler) Handle(ctx context.Context, command *v1.OrderPaidCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "orderPaidHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", command.GetAggregateID()))

	order, err := aggregate.LoadOrderAggregate(ctx, c.es, command.GetAggregateID())
	if err != nil {
		return err
	}

	if err := order.PayOrder(ctx, command); err != nil {
		return err
	}

	return c.es.Save(ctx, order)
}
