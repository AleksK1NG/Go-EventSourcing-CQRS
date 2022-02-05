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
		return a.onOrderCreated(evt)

	case events.OrderPaid:
		return a.onOrderPaid(evt)

	case events.OrderSubmitted:
		return a.onOrderSubmitted(evt)

	case events.OrderDelivered:
		return a.onOrderDelivered(evt)

	case events.OrderCanceled:
		return a.onOrderCanceled(evt)

	case events.OrderUpdated:
		return a.onOrderUpdated(evt)

	case events.OrderDeliveryAddressUpdated:
		return a.onOrderDeliveryAddressUpdated(evt)

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
		return a.onCreateOrderCommand(ctx, c)

	case *OrderPaidCommand:
		return a.onOrderPaidCommand(ctx, c)

	case *SubmitOrderCommand:
		return a.onSubmitOrderCommand(ctx, c)

	case *OrderUpdatedCommand:
		return a.onOrderUpdatedCommand(ctx, c)

	case *OrderCanceledCommand:
		return a.onOrderCanceledCommand(ctx, c)

	case *OrderDeliveredCommand:
		return a.onOrderDeliveredCommand(ctx, c)

	case *OrderChangeDeliveryAddressCommand:
		return a.onOrderChangeDeliveryAddressCommand(ctx, c)

	default:
		return es.ErrInvalidCommandType
	}
}
