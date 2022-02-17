package mappers

import (
	"github.com/AleksK1NG/es-microservice/internal/dto"
	"github.com/AleksK1NG/es-microservice/internal/order/events/v1"
)

func CreateOrderDtoToEventData(createDto dto.CreateOrderReqDto) v1.OrderCreatedEvent {
	return v1.OrderCreatedEvent{
		ShopItems:       createDto.ShopItems,
		AccountEmail:    createDto.AccountEmail,
		DeliveryAddress: createDto.DeliveryAddress,
	}
}
