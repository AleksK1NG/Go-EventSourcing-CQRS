package aggregate

import (
	"github.com/AleksK1NG/es-microservice/internal/order/events/v1"
	"github.com/AleksK1NG/es-microservice/internal/order/models"
	"github.com/AleksK1NG/es-microservice/pkg/es"
	"github.com/pkg/errors"
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
	aggregate.Order.ID = id
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

	case v1.OrderCreated:
		return a.onOrderCreated(evt)
	case v1.OrderPaid:
		return a.onOrderPaid(evt)
	case v1.OrderSubmitted:
		return a.onOrderSubmitted(evt)
	case v1.OrderCompleted:
		return a.onOrderCompleted(evt)
	case v1.OrderCanceled:
		return a.onOrderCanceled(evt)
	case v1.ShoppingCartUpdated:
		return a.onShoppingCartUpdated(evt)
	case v1.DeliveryAddressChanged:
		return a.onChangeDeliveryAddress(evt)

	default:
		return es.ErrInvalidEventType
	}
}

func (a *OrderAggregate) onOrderCreated(evt es.Event) error {
	var eventData v1.OrderCreatedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Order.AccountEmail = eventData.AccountEmail
	a.Order.ShopItems = eventData.ShopItems
	a.Order.TotalPrice = GetShopItemsTotalPrice(eventData.ShopItems)
	a.Order.DeliveryAddress = eventData.DeliveryAddress
	return nil
}

func (a *OrderAggregate) onOrderPaid(evt es.Event) error {
	var payment models.Payment
	if err := evt.GetJsonData(&payment); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Order.Paid = true
	a.Order.Payment = payment
	return nil
}

func (a *OrderAggregate) onOrderSubmitted(evt es.Event) error {
	a.Order.Submitted = true
	return nil
}

func (a *OrderAggregate) onOrderCompleted(evt es.Event) error {
	var eventData v1.OrderCompletedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Order.Completed = true
	a.Order.DeliveredTime = eventData.DeliveryTimestamp
	a.Order.Canceled = false
	return nil
}

func (a *OrderAggregate) onOrderCanceled(evt es.Event) error {
	var eventData v1.OrderCanceledEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Order.Canceled = true
	a.Order.Completed = false
	a.Order.CancelReason = eventData.CancelReason
	return nil
}

func (a *OrderAggregate) onShoppingCartUpdated(evt es.Event) error {
	var eventData v1.ShoppingCartUpdatedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Order.ShopItems = eventData.ShopItems
	a.Order.TotalPrice = GetShopItemsTotalPrice(eventData.ShopItems)
	return nil
}

func (a *OrderAggregate) onChangeDeliveryAddress(evt es.Event) error {
	var eventData v1.OrderDeliveryAddressChangedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Order.DeliveryAddress = eventData.DeliveryAddress
	return nil
}
