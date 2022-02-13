package mappers

import (
	"github.com/AleksK1NG/es-microservice/internal/dto"
	"github.com/AleksK1NG/es-microservice/internal/order/events"
)

func ChangeDeliveryAddressReqDtoToEventData(reqDto dto.ChangeDeliveryAddressReqDto) events.OrderChangeDeliveryAddress {
	return events.OrderChangeDeliveryAddress{
		DeliveryAddress: reqDto.DeliveryAddress,
	}
}
