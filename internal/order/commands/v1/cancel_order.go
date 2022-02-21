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

type CancelOrderCommandHandler interface {
	Handle(ctx context.Context, command *CancelOrderCommand) error
}

type cancelOrderCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  es.AggregateStore
}

func NewCancelOrderCommandHandler(log logger.Logger, cfg *config.Config, es es.AggregateStore) *cancelOrderCommandHandler {
	return &cancelOrderCommandHandler{log: log, cfg: cfg, es: es}
}

func (c *cancelOrderCommandHandler) Handle(ctx context.Context, command *CancelOrderCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cancelOrderCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", command.GetAggregateID()))

	order, err := aggregate.LoadOrderAggregate(ctx, c.es, command.GetAggregateID())
	if err != nil {
		return err
	}

	if err := order.CancelOrder(ctx, command.CancelReason); err != nil {
		return err
	}

	return c.es.Save(ctx, order)
}
