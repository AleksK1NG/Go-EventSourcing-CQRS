package aggregate

import "github.com/AleksK1NG/es-microservice/internal/order/events"

type CreateOrderCommand struct {
	events.OrderCreatedData
	AggregateID string `json:"aggregateID" validate:"required,gte=0"`
}

func NewCreateOrderCommand(orderCreatedData events.OrderCreatedData, aggregateID string) *CreateOrderCommand {
	return &CreateOrderCommand{OrderCreatedData: orderCreatedData, AggregateID: aggregateID}
}

func (o *CreateOrderCommand) GetAggregateID() string {
	return o.AggregateID
}

type OrderPaidCommand struct {
	AggregateID string `json:"aggregateID" validate:"required,gte=0"`
}

func NewOrderPaidCommand(aggregateID string) *OrderPaidCommand {
	return &OrderPaidCommand{AggregateID: aggregateID}
}

func (o *OrderPaidCommand) GetAggregateID() string {
	return o.AggregateID
}

type SubmitOrderCommand struct {
	AggregateID string `json:"aggregateID" validate:"required,gte=0"`
}

func NewSubmitOrderCommand(aggregateID string) *SubmitOrderCommand {
	return &SubmitOrderCommand{AggregateID: aggregateID}
}

func (o *SubmitOrderCommand) GetAggregateID() string {
	return o.AggregateID
}

type OrderUpdatedCommand struct {
	events.OrderUpdatedData
	AggregateID string `json:"aggregateID" validate:"required,gte=0"`
}

func NewOrderUpdatedCommand(orderUpdatedData events.OrderUpdatedData, aggregateID string) *OrderUpdatedCommand {
	return &OrderUpdatedCommand{OrderUpdatedData: orderUpdatedData, AggregateID: aggregateID}
}

func (o *OrderUpdatedCommand) GetAggregateID() string {
	return o.AggregateID
}
