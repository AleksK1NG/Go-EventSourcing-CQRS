package mappers

import (
	"github.com/AleksK1NG/es-microservice/internal/dto"
	orderService "github.com/AleksK1NG/es-microservice/proto/order"
)

func SearchResponseFromProto(protoSearch *orderService.SearchRes) dto.OrderSearchResponseDto {
	orders := make([]dto.OrderResponseDto, 0, len(protoSearch.GetOrders()))
	for _, order := range protoSearch.GetOrders() {
		orders = append(orders, OrderResponseDtoFromProto(order))
	}
	return dto.OrderSearchResponseDto{
		Pagination: PaginationFromProto(protoSearch.GetPagination()),
		Orders:     orders,
	}
}

func SearchResponseToProto(protoSearch *dto.OrderSearchResponseDto) *orderService.SearchRes {
	orders := make([]*orderService.Order, 0, len(protoSearch.Orders))
	for _, order := range protoSearch.Orders {
		orders = append(orders, OrderResponseDtoToProto(order))
	}
	return &orderService.SearchRes{
		Pagination: PaginationToProto(protoSearch.Pagination),
		Orders:     orders,
	}
}
