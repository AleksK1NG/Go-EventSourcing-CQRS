package aggregate

import (
	"github.com/AleksK1NG/es-microservice/internal/dto"
	"github.com/AleksK1NG/es-microservice/pkg/es"
	"github.com/pkg/errors"
)

func (a *OrderAggregate) onOrderCreated(evt es.Event) error {
	var eventData dto.OrderCreatedData
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

func (a *OrderAggregate) onOrderDelivering(evt es.Event) error {
	a.Order.Delivering = true
	return nil
}

func (a *OrderAggregate) onOrderDelivered(evt es.Event) error {
	a.Order.Delivered = true
	return nil
}

func (a *OrderAggregate) onOrderCanceled(evt es.Event) error {
	a.Order.Canceled = true
	a.Order.Delivered = false
	return nil
}

func (a *OrderAggregate) onOrderUpdated(evt es.Event) error {
	var eventData dto.OrderUpdatedData
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Order.ShopItems = eventData.ShopItems
	a.Order.TotalPrice = GetShopItemsTotalPrice(eventData.ShopItems)
	return nil
}
