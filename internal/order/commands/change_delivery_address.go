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

type ChangeOrderDeliveryAddressCommandHandler interface {
	Handle(ctx context.Context, command *v1.OrderChangeDeliveryAddressCommand) error
}

type changeOrderDeliveryAddressCmdHandler struct {
	log logger.Logger
	cfg *config.Config
	es  es.AggregateStore
}

func NewChangeOrderDeliveryAddressCmdHandler(log logger.Logger, cfg *config.Config, es es.AggregateStore) *changeOrderDeliveryAddressCmdHandler {
	return &changeOrderDeliveryAddressCmdHandler{log: log, cfg: cfg, es: es}
}

func (c *changeOrderDeliveryAddressCmdHandler) Handle(ctx context.Context, command *v1.OrderChangeDeliveryAddressCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "changeOrderDeliveryAddressCmdHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", command.GetAggregateID()))

	order, err := aggregate.LoadOrderAggregate(ctx, c.es, command.GetAggregateID())
	if err != nil {
		return err
	}

	if err := order.ChangeDeliveryAddress(ctx, command); err != nil {
		return err
	}

	return c.es.Save(ctx, order)
}
