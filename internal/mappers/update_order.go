package mappers

import (
	"github.com/AleksK1NG/es-microservice/internal/dto"
	"github.com/AleksK1NG/es-microservice/internal/order/events/v1"
)

func UpdateOrderReqDtoToEventData(reqDto dto.UpdateOrderItemsReqDto) v1.OrderUpdatedEvent {
	return v1.OrderUpdatedEvent{
		ShopItems: reqDto.ShopItems,
	}
}
