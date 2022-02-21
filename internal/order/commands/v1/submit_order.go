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

type SubmitOrderCommandHandler interface {
	Handle(ctx context.Context, command *SubmitOrderCommand) error
}

type submitOrderHandler struct {
	log logger.Logger
	cfg *config.Config
	es  es.AggregateStore
}

func NewSubmitOrderHandler(log logger.Logger, cfg *config.Config, es es.AggregateStore) *submitOrderHandler {
	return &submitOrderHandler{log: log, cfg: cfg, es: es}
}

func (c *submitOrderHandler) Handle(ctx context.Context, command *SubmitOrderCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "submitOrderHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", command.GetAggregateID()))

	order, err := aggregate.LoadOrderAggregate(ctx, c.es, command.GetAggregateID())
	if err != nil {
		return err
	}

	if err := order.SubmitOrder(ctx); err != nil {
		return err
	}

	return c.es.Save(ctx, order)
}
