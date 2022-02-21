package mappers

import (
	"github.com/AleksK1NG/es-microservice/internal/dto"
	"github.com/AleksK1NG/es-microservice/internal/order/events/v1"
)

func UpdateOrderReqDtoToEventData(reqDto dto.UpdateShoppingItemsReqDto) v1.ShoppingCartUpdatedEvent {
	return v1.ShoppingCartUpdatedEvent{
		ShopItems: reqDto.ShopItems,
	}
}
