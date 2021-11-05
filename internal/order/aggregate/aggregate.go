package aggregate

import (
	"github.com/AleksK1NG/es-microservice/internal/models"
	"github.com/AleksK1NG/es-microservice/internal/order/events"
	"github.com/AleksK1NG/es-microservice/pkg/es"
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

func (a *OrderAggregate) HandleCommand(command es.Command) error {

	switch c := command.(type) {

	case *CreateOrderCommand:
		return a.handleCreateOrderCommand(c)

	case *OrderPaidCommand:
		return a.handleOrderPaidCommand(c)

	case *SubmitOrderCommand:
		return a.handleSubmitOrderCommand(c)

	case *OrderUpdatedCommand:
		return a.handleOrderUpdatedCommand(c)

	default:
		return es.ErrInvalidCommandType
	}
}
