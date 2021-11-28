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
	span, ctx := opentracing.StartSpanFromContext(ctx, "submitOrderHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", command.GetAggregateID()))

	return aggregate.HandleCommandWithExists(ctx, c.es, command)
}
