package aggregate

import (
	"github.com/AleksK1NG/es-microservice/internal/order/events"
	"github.com/AleksK1NG/es-microservice/pkg/es"
)

func (a *OrderAggregate) handleOrderCreatedEvent(evt es.Event) error {
	var eventData events.OrderCreatedData
	if err := evt.GetJsonData(&eventData); err != nil {
		return err
	}

	a.Order.AccountEmail = eventData.AccountEmail
	a.Order.ShopItems = eventData.ShopItems
	a.Order.Created = true
	a.Order.TotalPrice = GetShopItemsTotalPrice(eventData.ShopItems)
	return nil
}

func (a *OrderAggregate) handleOrderPainEvent(evt es.Event) error {
	a.Order.Paid = true
	return nil
}

func (a *OrderAggregate) handleOrderSubmittedEvent(evt es.Event) error {
	a.Order.Submitted = true
	return nil
}

func (a *OrderAggregate) handleOrderDeliveringEvent(evt es.Event) error {
	a.Order.Delivering = true
	return nil
}

func (a *OrderAggregate) handleOrderDeliveredEvent(evt es.Event) error {
	a.Order.Delivered = true
	return nil
}

func (a *OrderAggregate) handleOrderCanceledEvent(evt es.Event) error {
	a.Order.Canceled = true
	a.Order.Delivered = false
	return nil
}

func (a *OrderAggregate) handleOrderUpdatedEvent(evt es.Event) error {
	var eventData events.OrderUpdatedData
	if err := evt.GetJsonData(&eventData); err != nil {
		return err
	}

	a.Order.ShopItems = eventData.ShopItems
	a.Order.TotalPrice = GetShopItemsTotalPrice(eventData.ShopItems)
	return nil
}
