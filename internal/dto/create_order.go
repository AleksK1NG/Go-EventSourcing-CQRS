package dto

import "github.com/AleksK1NG/es-microservice/internal/order/events"

type CreateOrderDto struct {
	events.OrderCreatedData
}
