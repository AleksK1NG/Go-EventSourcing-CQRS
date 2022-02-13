package mappers

import (
	"github.com/AleksK1NG/es-microservice/internal/dto"
	"github.com/AleksK1NG/es-microservice/internal/order/events"
)

func CreateOrderDtoToEventData(createDto dto.CreateOrderReqDto) events.OrderCreatedEventData {
	return events.OrderCreatedEventData{
		ShopItems:       createDto.ShopItems,
		AccountEmail:    createDto.AccountEmail,
		DeliveryAddress: createDto.DeliveryAddress,
	}
}
