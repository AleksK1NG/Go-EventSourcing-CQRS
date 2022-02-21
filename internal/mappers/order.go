package mappers

import (
	"github.com/AleksK1NG/es-microservice/internal/dto"
	"github.com/AleksK1NG/es-microservice/internal/order/aggregate"
	"github.com/AleksK1NG/es-microservice/internal/order/models"
	orderService "github.com/AleksK1NG/es-microservice/proto/order"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func OrderProjectionFromAggregate(orderAggregate *aggregate.OrderAggregate) *models.OrderProjection {
	return &models.OrderProjection{
		OrderID:         aggregate.GetOrderAggregateID(orderAggregate.GetID()),
		ShopItems:       orderAggregate.Order.ShopItems,
		Paid:            orderAggregate.Order.Paid,
		Submitted:       orderAggregate.Order.Submitted,
		Completed:       orderAggregate.Order.Completed,
		Canceled:        orderAggregate.Order.Canceled,
		AccountEmail:    orderAggregate.Order.AccountEmail,
		TotalPrice:      orderAggregate.Order.TotalPrice,
		DeliveredTime:   orderAggregate.Order.DeliveredTime,
		CancelReason:    orderAggregate.Order.CancelReason,
		DeliveryAddress: orderAggregate.Order.DeliveryAddress,
		Payment:         orderAggregate.Order.Payment,
	}
}

func OrderResponseFromProjection(projection *models.OrderProjection) dto.OrderResponseDto {
	return dto.OrderResponseDto{
		ID:              projection.ID,
		OrderID:         projection.OrderID,
		ShopItems:       ShopItemsResponseFromModels(projection.ShopItems),
		AccountEmail:    projection.AccountEmail,
		DeliveryAddress: projection.DeliveryAddress,
		CancelReason:    projection.CancelReason,
		TotalPrice:      projection.TotalPrice,
		DeliveredTime:   projection.DeliveredTime,
		Paid:            projection.Paid,
		Submitted:       projection.Submitted,
		Completed:       projection.Completed,
		Canceled:        projection.Canceled,
		Payment:         PaymentResponseFromModel(projection.Payment),
	}
}

func OrderResponseDtoFromProto(orderProto *orderService.Order) dto.OrderResponseDto {
	return dto.OrderResponseDto{
		OrderID:         orderProto.GetID(),
		ShopItems:       ShopItemsResponseFromProto(orderProto.GetShopItems()),
		AccountEmail:    orderProto.GetAccountEmail(),
		DeliveryAddress: orderProto.GetDeliveryAddress(),
		CancelReason:    orderProto.GetCancelReason(),
		TotalPrice:      orderProto.GetTotalPrice(),
		DeliveredTime:   orderProto.GetDeliveryTimestamp().AsTime(),
		Paid:            orderProto.GetPaid(),
		Submitted:       orderProto.GetSubmitted(),
		Completed:       orderProto.GetCompleted(),
		Canceled:        orderProto.GetCanceled(),
		Payment:         PaymentFromProto(orderProto.GetPayment()),
	}
}

func OrdersFromProjections(projections []*models.OrderProjection) []dto.OrderResponseDto {
	orders := make([]dto.OrderResponseDto, 0, len(projections))
	for _, projection := range projections {
		orders = append(orders, OrderResponseFromProjection(projection))
	}
	return orders
}

func OrderResponseDtoToProto(orderDto dto.OrderResponseDto) *orderService.Order {
	return &orderService.Order{
		ID:                orderDto.OrderID,
		ShopItems:         ShopItemsResponseToProto(orderDto.ShopItems),
		Paid:              orderDto.Paid,
		Submitted:         orderDto.Submitted,
		Completed:         orderDto.Completed,
		Canceled:          orderDto.Canceled,
		TotalPrice:        orderDto.TotalPrice,
		AccountEmail:      orderDto.AccountEmail,
		CancelReason:      orderDto.CancelReason,
		DeliveryAddress:   orderDto.DeliveryAddress,
		DeliveryTimestamp: timestamppb.New(orderDto.DeliveredTime),
		Payment:           PaymentToProto(orderDto.Payment),
	}
}

func OrdersResponseDtoToProto(ordersDto []dto.OrderResponseDto) []*orderService.Order {
	orders := make([]*orderService.Order, 0, len(ordersDto))
	for _, order := range ordersDto {
		orders = append(orders, OrderResponseDtoToProto(order))
	}
	return orders
}
