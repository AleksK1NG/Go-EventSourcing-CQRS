package aggregate

import (
	"context"
	"github.com/AleksK1NG/es-microservice/internal/models"
	"github.com/AleksK1NG/es-microservice/internal/order/events"
	"github.com/AleksK1NG/es-microservice/pkg/es"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

const (
	OrderAggregateType es.AggregateType = "order"
)

type OrderAggregateInterface interface {
	HandleCreateEvent(createEvent es.Event) error
	HandlePaidEvent(payEvent es.Event) error
	HandleSubmitEvent(submitEvent es.Event) error
}

type OrderAggregate struct {
	*es.AggregateBase
	Order *models.Order
}

func NewOrderAggregateWithID(id string) *OrderAggregate {
	if id == "" {
		return nil
	}

	aggregate := NewOrderAggregate()
	aggregate.SetID(id)
	return aggregate
}

func NewOrderAggregate() *OrderAggregate {
	orderAggregate := &OrderAggregate{Order: models.NewOrder()}
	base := es.NewAggregateBase(orderAggregate.When)
	base.SetType(OrderAggregateType)
	orderAggregate.AggregateBase = base
	return orderAggregate
}

func (a *OrderAggregate) When(evt es.Event) error {

	switch evt.GetEventType() {

	case events.OrderCreated:
		return a.handleOrderCreatedEvent(evt)

	case events.OrderPaid:
		return a.handleOrderPainEvent(evt)

	case events.OrderSubmitted:
		return a.handleOrderSubmittedEvent(evt)

	case events.OrderDelivering:
		return a.handleOrderDeliveringEvent(evt)

	case events.OrderDelivered:
		return a.handleOrderDeliveredEvent(evt)

	case events.OrderCanceled:
		return a.handleOrderCanceledEvent(evt)

	case events.OrderUpdated:
		return a.handleOrderUpdatedEvent(evt)

	default:
		return es.ErrInvalidEventType
	}
}

func (a *OrderAggregate) HandleCommand(ctx context.Context, command es.Command) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderAggregate.HandleCommand")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", command.GetAggregateID()))

	switch c := command.(type) {

	case *CreateOrderCommand:
		return a.handleCreateOrderCommand(ctx, c)

	case *OrderPaidCommand:
		return a.handleOrderPaidCommand(ctx, c)

	case *SubmitOrderCommand:
		return a.handleSubmitOrderCommand(ctx, c)

	case *OrderUpdatedCommand:
		return a.handleOrderUpdatedCommand(ctx, c)

	default:
		return es.ErrInvalidCommandType
	}
}
