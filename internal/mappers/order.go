package mappers

import (
	"github.com/AleksK1NG/es-microservice/internal/models"
	"github.com/AleksK1NG/es-microservice/internal/order/aggregate"
)

func OrderProjectionFromAggregate(orderAggregate *aggregate.OrderAggregate) *models.OrderProjection {
	return &models.OrderProjection{
		OrderID:    orderAggregate.GetID(),
		ItemsIDs:   orderAggregate.Order.ItemsIDs,
		Created:    orderAggregate.Order.Created,
		Paid:       orderAggregate.Order.Paid,
		Submitted:  orderAggregate.Order.Submitted,
		Delivering: orderAggregate.Order.Delivering,
		Delivered:  orderAggregate.Order.Delivered,
		Canceled:   orderAggregate.Order.Canceled,
	}
}
