package mappers

import (
	"github.com/AleksK1NG/es-microservice/internal/dto"
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
		Payment:         orderAggregate.Order.Payment,
	}
}

func GetOrderResponseFromProjection(projection *models.OrderProjection) dto.GetOrderResponseDto {
	return dto.GetOrderResponseDto{
		ID:              projection.ID,
		OrderID:         projection.OrderID,
		ShopItems:       projection.ShopItems,
		AccountEmail:    projection.AccountEmail,
		DeliveryAddress: projection.DeliveryAddress,
		CancelReason:    projection.CancelReason,
		TotalPrice:      projection.TotalPrice,
		DeliveredTime:   projection.DeliveredTime,
		Created:         projection.Created,
		Paid:            projection.Paid,
		Submitted:       projection.Submitted,
		Delivered:       projection.Delivered,
		Canceled:        projection.Canceled,
		Payment:         projection.Payment,
	}
}
