package aggregate

import (
	"github.com/AleksK1NG/es-microservice/internal/models"
	"github.com/AleksK1NG/es-microservice/internal/order/events"
	"github.com/AleksK1NG/es-microservice/pkg/es"
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
		return a.onOrderChangeDeliveryAddress(evt)

	default:
		return es.ErrInvalidEventType
	}
}

//func (a *OrderAggregate) HandleCommand(ctx context.Context, command es.Command) error {
//	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderAggregate.HandleCommand")
//	defer span.Finish()
//	span.LogFields(log.String("AggregateID", command.GetAggregateID()))
//
//	switch c := command.(type) {
//
//	case *CreateOrderCommand:
//		return a.CreateOrder(ctx, c)
//
//	case *OrderPaidCommand:
//		return a.PayOrder(ctx, c)
//
//	case *SubmitOrderCommand:
//		return a.SubmitOrder(ctx, c)
//
//	case *OrderUpdatedCommand:
//		return a.UpdateOrder(ctx, c)
//
//	case *OrderCanceledCommand:
//		return a.CancelOrder(ctx, c)
//
//	case *OrderDeliveredCommand:
//		return a.DeliverOrder(ctx, c)
//
//	case *OrderChangeDeliveryAddressCommand:
//		return a.ChangeDeliveryAddress(ctx, c)
//
//	default:
//		return es.ErrInvalidCommandType
//	}
//}
