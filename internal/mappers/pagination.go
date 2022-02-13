package mappers

import (
	"github.com/AleksK1NG/es-microservice/internal/dto"
	orderService "github.com/AleksK1NG/es-microservice/proto/order"
)

func PaginationFromProto(protoPagination *orderService.Pagination) dto.Pagination {
	return dto.Pagination{
		TotalCount: protoPagination.GetTotalCount(),
		TotalPages: protoPagination.GetTotalPages(),
		Page:       protoPagination.GetPage(),
		Size:       protoPagination.GetSize(),
		HasMore:    protoPagination.GetHasMore(),
	}
}

func PaginationToProto(protoPagination dto.Pagination) *orderService.Pagination {
	return &orderService.Pagination{
		TotalCount: protoPagination.TotalCount,
		TotalPages: protoPagination.TotalPages,
		Page:       protoPagination.Page,
		Size:       protoPagination.Size,
		HasMore:    protoPagination.HasMore,
	}
}
