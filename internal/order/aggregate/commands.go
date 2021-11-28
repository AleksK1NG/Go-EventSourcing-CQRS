package aggregate

import (
	"github.com/AleksK1NG/es-microservice/internal/dto"
	"github.com/AleksK1NG/es-microservice/pkg/es"
)

type CreateOrderCommand struct {
	dto.OrderCreatedData
	es.BaseCommand
}

func NewCreateOrderCommand(orderCreatedData dto.OrderCreatedData, aggregateID string) *CreateOrderCommand {
	return &CreateOrderCommand{OrderCreatedData: orderCreatedData, BaseCommand: es.NewBaseCommand(aggregateID)}
}

type OrderPaidCommand struct {
	es.BaseCommand
}

func NewOrderPaidCommand(aggregateID string) *OrderPaidCommand {
	return &OrderPaidCommand{BaseCommand: es.NewBaseCommand(aggregateID)}
}

type SubmitOrderCommand struct {
	es.BaseCommand
}

func NewSubmitOrderCommand(aggregateID string) *SubmitOrderCommand {
	return &SubmitOrderCommand{BaseCommand: es.NewBaseCommand(aggregateID)}
}

type OrderUpdatedCommand struct {
	dto.OrderUpdatedData
	es.BaseCommand
}

func NewOrderUpdatedCommand(orderUpdatedData dto.OrderUpdatedData, aggregateID string) *OrderUpdatedCommand {
	return &OrderUpdatedCommand{OrderUpdatedData: orderUpdatedData, BaseCommand: es.NewBaseCommand(aggregateID)}
}
