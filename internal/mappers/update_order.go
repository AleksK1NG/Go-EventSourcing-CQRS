package mappers

import (
	"github.com/AleksK1NG/es-microservice/internal/dto"
	"github.com/AleksK1NG/es-microservice/internal/order/events/v1"
)

func UpdateOrderReqDtoToEventData(reqDto dto.UpdateOrderItemsReqDto) v1.OrderUpdatedEventData {
	return v1.OrderUpdatedEventData{
		ShopItems: reqDto.ShopItems,
	}
}
