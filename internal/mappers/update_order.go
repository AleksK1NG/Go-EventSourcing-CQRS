package mappers

import (
	"github.com/AleksK1NG/es-microservice/internal/dto"
	"github.com/AleksK1NG/es-microservice/internal/order/events"
)

func UpdateOrderReqDtoToEventData(reqDto dto.UpdateOrderItemsReqDto) events.OrderUpdatedEventData {
	return events.OrderUpdatedEventData{
		ShopItems: reqDto.ShopItems,
	}
}
