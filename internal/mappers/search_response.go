package mappers

import (
	"github.com/AleksK1NG/es-microservice/internal/dto"
	orderService "github.com/AleksK1NG/es-microservice/proto/order"
)

func SearchResponseFromProto(protoSearch *orderService.SearchRes) dto.OrderSearchResponseDto {
	orders := make([]dto.GetOrderResponseDto, 0, len(protoSearch.GetOrders()))
	for _, order := range protoSearch.GetOrders() {
		orders = append(orders, OrderResponseDtoFromProto(order))
	}
	return dto.OrderSearchResponseDto{
		Pagination: PaginationFromProto(protoSearch.GetPagination()),
		Orders:     orders,
	}
}
