package events

import (
	"github.com/AleksK1NG/es-microservice/pkg/es"
)

const (
	OrderCreated                = "V1_ORDER_CREATED"
	OrderPaid                   = "V1_ORDER_PAID"
	OrderSubmitted              = "V1_ORDER_SUBMITTED"
	OrderDelivering             = "V1_ORDER_DELIVERING"
	OrderDelivered              = "V1_ORDER_DELIVERED"
	OrderCanceled               = "V1_ORDER_CANCELED"
	OrderUpdated                = "V1_ORDER_UPDATED"
	OrderDeliveryAddressUpdated = "V1_ORDER_DELIVERY_ADDRESS_UPDATED"
)

func NewCreateOrderEvent(aggregate es.Aggregate, data []byte) es.Event {
	createOrderEvent := es.NewBaseEvent(aggregate, OrderCreated)
	createOrderEvent.SetData(data)
	return createOrderEvent
}

func NewPayOrderEvent(aggregate es.Aggregate, data []byte) es.Event {
	orderPaidEvent := es.NewBaseEvent(aggregate, OrderPaid)
	orderPaidEvent.SetData(data)
	return orderPaidEvent
}

func NewSubmitOrderEvent(aggregate es.Aggregate) es.Event {
	return es.NewBaseEvent(aggregate, OrderSubmitted)
}

func NewOrderUpdatedEvent(aggregate es.Aggregate, data []byte) es.Event {
	orderUpdatedEvent := es.NewBaseEvent(aggregate, OrderUpdated)
	orderUpdatedEvent.SetData(data)
	return orderUpdatedEvent
}

func NewOrderDeliveryAddressUpdatedEvent(aggregate es.Aggregate, data []byte) es.Event {
	orderUpdatedEvent := es.NewBaseEvent(aggregate, OrderDeliveryAddressUpdated)
	orderUpdatedEvent.SetData(data)
	return orderUpdatedEvent
}

func NewOrderCanceledEvent(aggregate es.Aggregate, data []byte) es.Event {
	orderUpdatedEvent := es.NewBaseEvent(aggregate, OrderCanceled)
	orderUpdatedEvent.SetData(data)
	return orderUpdatedEvent
}

func NewOrderDeliveredEvent(aggregate es.Aggregate, data []byte) es.Event {
	orderUpdatedEvent := es.NewBaseEvent(aggregate, OrderDelivered)
	orderUpdatedEvent.SetData(data)
	return orderUpdatedEvent
}
