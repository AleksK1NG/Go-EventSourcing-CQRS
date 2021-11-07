package dto

import "github.com/AleksK1NG/es-microservice/internal/order/events"

type UpdateOrderDto struct {
	events.OrderUpdatedData
}
