package v1

import (
	"github.com/AleksK1NG/es-microservice/internal/order/events/v1"
	"github.com/AleksK1NG/es-microservice/internal/order/models"
	"github.com/AleksK1NG/es-microservice/pkg/es"
)

type CreateOrderCommand struct {
	v1.OrderCreatedEvent
	es.BaseCommand
}

func NewCreateOrderCommand(orderCreatedData v1.OrderCreatedEvent, aggregateID string) *CreateOrderCommand {
	return &CreateOrderCommand{OrderCreatedEvent: orderCreatedData, BaseCommand: es.NewBaseCommand(aggregateID)}
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
	v1.OrderUpdatedEvent
	es.BaseCommand
}

func NewOrderUpdatedCommand(orderUpdatedData v1.OrderUpdatedEvent, aggregateID string) *OrderUpdatedCommand {
	return &OrderUpdatedCommand{OrderUpdatedEvent: orderUpdatedData, BaseCommand: es.NewBaseCommand(aggregateID)}
}

type OrderCanceledCommand struct {
	v1.OrderCanceledEvent
	es.BaseCommand
}

func NewOrderCanceledCommand(orderCanceledEventData v1.OrderCanceledEvent, aggregateID string) *OrderCanceledCommand {
	return &OrderCanceledCommand{OrderCanceledEvent: orderCanceledEventData, BaseCommand: es.NewBaseCommand(aggregateID)}
}

type OrderDeliveredCommand struct {
	v1.OrderDeliveredEvent
	es.BaseCommand
}

func NewOrderDeliveredCommand(orderDeliveredEventData v1.OrderDeliveredEvent, aggregateID string) *OrderDeliveredCommand {
	return &OrderDeliveredCommand{OrderDeliveredEvent: orderDeliveredEventData, BaseCommand: es.NewBaseCommand(aggregateID)}
}

type OrderChangeDeliveryAddressCommand struct {
	v1.OrderDeliveryAddressChangedEvent
	es.BaseCommand
}

func NewOrderChangeDeliveryAddressCommand(orderChangeDeliveryAddress v1.OrderDeliveryAddressChangedEvent, aggregateID string) *OrderChangeDeliveryAddressCommand {
	return &OrderChangeDeliveryAddressCommand{OrderDeliveryAddressChangedEvent: orderChangeDeliveryAddress, BaseCommand: es.NewBaseCommand(aggregateID)}
}
