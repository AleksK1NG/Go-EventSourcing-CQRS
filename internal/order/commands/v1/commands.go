package v1

import (
	"github.com/AleksK1NG/es-microservice/internal/order/events/v1"
	"github.com/AleksK1NG/es-microservice/internal/order/models"
	"github.com/AleksK1NG/es-microservice/pkg/es"
)

type CreateOrderCommand struct {
	v1.OrderCreatedEventData
	es.BaseCommand
}

func NewCreateOrderCommand(orderCreatedData v1.OrderCreatedEventData, aggregateID string) *CreateOrderCommand {
	return &CreateOrderCommand{OrderCreatedEventData: orderCreatedData, BaseCommand: es.NewBaseCommand(aggregateID)}
}

type OrderPaidCommand struct {
	models.Payment
	es.BaseCommand
}

func NewOrderPaidCommand(payment models.Payment, aggregateID string) *OrderPaidCommand {
	return &OrderPaidCommand{Payment: payment, BaseCommand: es.NewBaseCommand(aggregateID)}
}

type SubmitOrderCommand struct {
	es.BaseCommand
}

func NewSubmitOrderCommand(aggregateID string) *SubmitOrderCommand {
	return &SubmitOrderCommand{BaseCommand: es.NewBaseCommand(aggregateID)}
}

type OrderUpdatedCommand struct {
	v1.OrderUpdatedEventData
	es.BaseCommand
}

func NewOrderUpdatedCommand(orderUpdatedData v1.OrderUpdatedEventData, aggregateID string) *OrderUpdatedCommand {
	return &OrderUpdatedCommand{OrderUpdatedEventData: orderUpdatedData, BaseCommand: es.NewBaseCommand(aggregateID)}
}

type OrderCanceledCommand struct {
	v1.OrderCanceledEventData
	es.BaseCommand
}

func NewOrderCanceledCommand(orderCanceledEventData v1.OrderCanceledEventData, aggregateID string) *OrderCanceledCommand {
	return &OrderCanceledCommand{OrderCanceledEventData: orderCanceledEventData, BaseCommand: es.NewBaseCommand(aggregateID)}
}

type OrderDeliveredCommand struct {
	v1.OrderDeliveredEventData
	es.BaseCommand
}

func NewOrderDeliveredCommand(orderDeliveredEventData v1.OrderDeliveredEventData, aggregateID string) *OrderDeliveredCommand {
	return &OrderDeliveredCommand{OrderDeliveredEventData: orderDeliveredEventData, BaseCommand: es.NewBaseCommand(aggregateID)}
}

type OrderChangeDeliveryAddressCommand struct {
	v1.OrderChangeDeliveryAddress
	es.BaseCommand
}

func NewOrderChangeDeliveryAddressCommand(orderChangeDeliveryAddress v1.OrderChangeDeliveryAddress, aggregateID string) *OrderChangeDeliveryAddressCommand {
	return &OrderChangeDeliveryAddressCommand{OrderChangeDeliveryAddress: orderChangeDeliveryAddress, BaseCommand: es.NewBaseCommand(aggregateID)}
}
