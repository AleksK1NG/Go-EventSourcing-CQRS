package commands

import (
	"context"
	"github.com/AleksK1NG/es-microservice/config"
	"github.com/AleksK1NG/es-microservice/internal/order/aggregate"
	"github.com/AleksK1NG/es-microservice/pkg/es"
	"github.com/AleksK1NG/es-microservice/pkg/logger"
)

type OrderPaidCommandHandler interface {
	Handle(ctx context.Context, command *aggregate.OrderPaidCommand) error
}

type orderPaidHandler struct {
	log logger.Logger
	cfg *config.Config
	es  es.AggregateStore
}

func NewOrderPaidHandler(log logger.Logger, cfg *config.Config, es es.AggregateStore) *orderPaidHandler {
	return &orderPaidHandler{log: log, cfg: cfg, es: es}
}

func (c *orderPaidHandler) Handle(ctx context.Context, command *aggregate.OrderPaidCommand) error {
	order := aggregate.NewOrderAggregateWithID(command.AggregateID)
	err := c.es.Exists(ctx, order.GetID())
	if err != nil {
		return err
	}

	if err := c.es.Load(ctx, order); err != nil {
		return err
	}

	if err := order.HandleCommand(command); err != nil {
		return err
	}

	return c.es.Save(ctx, order)
}
