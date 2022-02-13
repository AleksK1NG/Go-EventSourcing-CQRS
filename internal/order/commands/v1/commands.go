package v1

import (
	"github.com/AleksK1NG/es-microservice/internal/order/events"
	"github.com/AleksK1NG/es-microservice/internal/order/models"
	"github.com/AleksK1NG/es-microservice/pkg/es"
)

type CreateOrderCommand struct {
	events.OrderCreatedEventData
	es.BaseCommand
}

func NewCreateOrderCommand(orderCreatedData events.OrderCreatedEventData, aggregateID string) *CreateOrderCommand {
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
	events.OrderUpdatedEventData
	es.BaseCommand
}

func NewOrderUpdatedCommand(orderUpdatedData events.OrderUpdatedEventData, aggregateID string) *OrderUpdatedCommand {
	return &OrderUpdatedCommand{OrderUpdatedEventData: orderUpdatedData, BaseCommand: es.NewBaseCommand(aggregateID)}
}

type OrderCanceledCommand struct {
	events.OrderCanceledEventData
	es.BaseCommand
}

func NewOrderCanceledCommand(orderCanceledEventData events.OrderCanceledEventData, aggregateID string) *OrderCanceledCommand {
	return &OrderCanceledCommand{OrderCanceledEventData: orderCanceledEventData, BaseCommand: es.NewBaseCommand(aggregateID)}
}

type OrderDeliveredCommand struct {
	events.OrderDeliveredEventData
	es.BaseCommand
}

func NewOrderDeliveredCommand(orderDeliveredEventData events.OrderDeliveredEventData, aggregateID string) *OrderDeliveredCommand {
	return &OrderDeliveredCommand{OrderDeliveredEventData: orderDeliveredEventData, BaseCommand: es.NewBaseCommand(aggregateID)}
}

type OrderChangeDeliveryAddressCommand struct {
	events.OrderChangeDeliveryAddress
	es.BaseCommand
}

func NewOrderChangeDeliveryAddressCommand(orderChangeDeliveryAddress events.OrderChangeDeliveryAddress, aggregateID string) *OrderChangeDeliveryAddressCommand {
	return &OrderChangeDeliveryAddressCommand{OrderChangeDeliveryAddress: orderChangeDeliveryAddress, BaseCommand: es.NewBaseCommand(aggregateID)}
}
