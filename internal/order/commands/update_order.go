package commands

import (
	"context"
	"github.com/AleksK1NG/es-microservice/config"
	"github.com/AleksK1NG/es-microservice/internal/order/aggregate"
	"github.com/AleksK1NG/es-microservice/pkg/es"
	"github.com/AleksK1NG/es-microservice/pkg/logger"
	"github.com/AleksK1NG/es-microservice/pkg/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type UpdateOrderCommandHandler interface {
	Handle(ctx context.Context, command *aggregate.OrderUpdatedCommand) error
}

type updateOrderCmdHandler struct {
	log logger.Logger
	cfg *config.Config
	es  es.AggregateStore
}

func NewUpdateOrderCmdHandler(log logger.Logger, cfg *config.Config, es es.AggregateStore) *updateOrderCmdHandler {
	return &updateOrderCmdHandler{log: log, cfg: cfg, es: es}
}

func (c *updateOrderCmdHandler) Handle(ctx context.Context, command *aggregate.OrderUpdatedCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "updateOrderCmdHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", command.GetAggregateID()))

	c.log.Infof("(TEXT MAP CARRIER): %+v", tracing.ExtractTextMapCarrier(span.Context()))

	order := aggregate.NewOrderAggregateWithID(command.AggregateID)
	err := c.es.Exists(ctx, order.GetID())
	if err != nil {
		return err
	}

	if err := c.es.Load(ctx, order); err != nil {
		return err
	}

	if err := order.HandleCommand(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return c.es.Save(ctx, order)
}
