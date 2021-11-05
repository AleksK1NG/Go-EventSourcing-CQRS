package commands

import (
	"context"
	"github.com/AleksK1NG/es-microservice/config"
	"github.com/AleksK1NG/es-microservice/internal/order/aggregate"
	"github.com/AleksK1NG/es-microservice/pkg/es"
	"github.com/AleksK1NG/es-microservice/pkg/logger"
)

type SubmitOrderCommandHandler interface {
	Handle(ctx context.Context, command *aggregate.SubmitOrderCommand) error
}

type submitOrderHandler struct {
	log logger.Logger
	cfg *config.Config
	es  es.AggregateStore
}

func NewSubmitOrderHandler(log logger.Logger, cfg *config.Config, es es.AggregateStore) *submitOrderHandler {
	return &submitOrderHandler{log: log, cfg: cfg, es: es}
}

func (c *submitOrderHandler) Handle(ctx context.Context, command *aggregate.SubmitOrderCommand) error {
	err := c.es.Exists(ctx, command.AggregateID)
	if err != nil {
		return err
	}

	order := aggregate.NewOrderAggregateWithID(command.AggregateID)

	if err := c.es.Load(ctx, order.GetID(), order); err != nil {
		return err
	}

	if err := order.HandleCommand(command); err != nil {
		return err
	}

	return c.es.Save(ctx, order)
}
