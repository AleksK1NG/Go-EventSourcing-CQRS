package projection

import (
	"context"
	"github.com/AleksK1NG/es-microservice/internal/models"
	"github.com/AleksK1NG/es-microservice/internal/order/events"
	"github.com/AleksK1NG/es-microservice/pkg/es"
)

func (o *orderProjection) handleOrderCreateEvent(ctx context.Context, evt es.Event) error {
	var eventData events.OrderCreatedData
	if err := evt.GetJsonData(&eventData); err != nil {
		return err
	}

	op := &models.OrderProjection{
		OrderID:      GetOrderAggregateID(evt.AggregateID),
		ShopItems:    eventData.ShopItems,
		Created:      true,
		Paid:         false,
		Submitted:    false,
		Delivering:   false,
		Delivered:    false,
		Canceled:     false,
		AccountEmail: eventData.AccountEmail,
		TotalPrice:   GetShopItemsTotalPrice(eventData.ShopItems),
	}

	result, err := o.mongoRepo.Insert(ctx, op)
	if err != nil {
		return err
	}

	o.log.Debugf("projection OrderCreated result: %s", result)
	return nil
}

func (o *orderProjection) handleOrderPaidEvent(ctx context.Context, evt es.Event) error {
	op := &models.OrderProjection{OrderID: GetOrderAggregateID(evt.AggregateID), Paid: true}
	return o.mongoRepo.UpdateOrder(ctx, op)
}

func (o *orderProjection) handleSubmitEvent(ctx context.Context, evt es.Event) error {
	op := &models.OrderProjection{OrderID: GetOrderAggregateID(evt.AggregateID), Submitted: true}
	return o.mongoRepo.UpdateOrder(ctx, op)
}

func (o *orderProjection) handleUpdateEvent(ctx context.Context, evt es.Event) error {
	var eventData events.OrderUpdatedData
	if err := evt.GetJsonData(&eventData); err != nil {
		return err
	}

	op := &models.OrderProjection{OrderID: GetOrderAggregateID(evt.AggregateID), ShopItems: eventData.ShopItems}
	op.TotalPrice = GetShopItemsTotalPrice(eventData.ShopItems)
	return o.mongoRepo.UpdateOrder(ctx, op)
}

func GetShopItemsTotalPrice(shopItems []*models.ShopItem) float64 {
	var totalPrice float64 = 0
	for _, item := range shopItems {
		totalPrice += item.Price * float64(item.Quantity)
	}
	return totalPrice
}
