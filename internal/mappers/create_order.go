package mappers

import (
	"github.com/AleksK1NG/es-microservice/internal/dto"
	"github.com/AleksK1NG/es-microservice/internal/order/events/v1"
)

func CreateOrderDtoToEventData(createDto dto.CreateOrderReqDto) v1.OrderCreatedEventData {
	return v1.OrderCreatedEventData{
		ShopItems:       createDto.ShopItems,
		AccountEmail:    createDto.AccountEmail,
		DeliveryAddress: createDto.DeliveryAddress,
	}
}
