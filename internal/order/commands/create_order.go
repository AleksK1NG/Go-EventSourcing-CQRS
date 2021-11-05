package commands

import (
	"context"
	"github.com/AleksK1NG/es-microservice/config"
	"github.com/AleksK1NG/es-microservice/internal/order/aggregate"
	"github.com/AleksK1NG/es-microservice/pkg/es"
	"github.com/AleksK1NG/es-microservice/pkg/logger"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/pkg/errors"
)

type CreateOrderCommandHandler interface {
	Handle(ctx context.Context, command *aggregate.CreateOrderCommand) error
}

type createOrderHandler struct {
	log logger.Logger
	cfg *config.Config
	es  es.AggregateStore
}

func NewCreateOrderHandler(log logger.Logger, cfg *config.Config, es es.AggregateStore) *createOrderHandler {
	return &createOrderHandler{log: log, cfg: cfg, es: es}
}

func (c *createOrderHandler) Handle(ctx context.Context, command *aggregate.CreateOrderCommand) error {
	err := c.es.Exists(ctx, command.AggregateID)
	if err != nil && !errors.Is(err, esdb.ErrStreamNotFound) {
		return err
	}

	order := aggregate.NewOrderAggregateWithID(command.AggregateID)

	if err := order.HandleCommand(command); err != nil {
		return err
	}

	return c.es.Save(ctx, order)
}
