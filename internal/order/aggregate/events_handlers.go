package aggregate

import (
	"github.com/AleksK1NG/es-microservice/internal/order/events"
	"github.com/AleksK1NG/es-microservice/pkg/es"
	"github.com/pkg/errors"
)

func (a *OrderAggregate) onOrderCreated(evt es.Event) error {
	var eventData events.OrderCreatedEventData
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Order.AccountEmail = eventData.AccountEmail
	a.Order.ShopItems = eventData.ShopItems
	a.Order.Created = true
	a.Order.TotalPrice = GetShopItemsTotalPrice(eventData.ShopItems)
	return nil
}

func (a *OrderAggregate) onOrderPaid(evt es.Event) error {
	a.Order.Paid = true
	return nil
}

func (a *OrderAggregate) onOrderSubmitted(evt es.Event) error {
	a.Order.Submitted = true
	return nil
}

func (a *OrderAggregate) onOrderDelivered(evt es.Event) error {
	var eventData events.OrderDeliveredEventData
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Order.Delivered = true
	a.Order.DeliveredTime = eventData.DeliveryTimestamp
	a.Order.Canceled = false
	return nil
}

func (a *OrderAggregate) onOrderCanceled(evt es.Event) error {
	var eventData events.OrderCanceledEventData
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Order.Canceled = true
	a.Order.Delivered = false
	a.Order.CancelReason = eventData.CancelReason
	return nil
}

func (a *OrderAggregate) onOrderUpdated(evt es.Event) error {
	var eventData events.OrderUpdatedEventData
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Order.ShopItems = eventData.ShopItems
	a.Order.TotalPrice = GetShopItemsTotalPrice(eventData.ShopItems)
	return nil
}
