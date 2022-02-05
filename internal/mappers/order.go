package mappers

import (
	"github.com/AleksK1NG/es-microservice/internal/order/aggregate"
	"github.com/AleksK1NG/es-microservice/internal/order/models"
)

func OrderProjectionFromAggregate(orderAggregate *aggregate.OrderAggregate) *models.OrderProjection {
	return &models.OrderProjection{
		OrderID:         aggregate.GetOrderAggregateID(orderAggregate.GetID()),
		ShopItems:       orderAggregate.Order.ShopItems,
		Created:         orderAggregate.Order.Created,
		Paid:            orderAggregate.Order.Paid,
		Submitted:       orderAggregate.Order.Submitted,
		Delivered:       orderAggregate.Order.Delivered,
		Canceled:        orderAggregate.Order.Canceled,
		AccountEmail:    orderAggregate.Order.AccountEmail,
		TotalPrice:      orderAggregate.Order.TotalPrice,
		DeliveredTime:   orderAggregate.Order.DeliveredTime,
		CancelReason:    orderAggregate.Order.CancelReason,
		DeliveryAddress: orderAggregate.Order.DeliveryAddress,
	}
}
