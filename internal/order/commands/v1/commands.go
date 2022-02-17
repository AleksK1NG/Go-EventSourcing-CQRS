package v1

import (
	"time"

	"github.com/AleksK1NG/es-microservice/internal/order/models"
	"github.com/AleksK1NG/es-microservice/pkg/es"
)

type CreateOrderCommand struct {
	es.BaseCommand
	ShopItems       []*models.ShopItem `json:"shopItems" bson:"shopItems,omitempty" validate:"required"`
	AccountEmail    string             `json:"accountEmail" bson:"accountEmail,omitempty" validate:"required,email"`
	DeliveryAddress string             `json:"deliveryAddress" bson:"deliveryAddress,omitempty" validate:"required"`
}

func NewCreateOrderCommand(aggregateID string, shopItems []*models.ShopItem, accountEmail, deliveryAddress string) *CreateOrderCommand {
	return &CreateOrderCommand{BaseCommand: es.NewBaseCommand(aggregateID), ShopItems: shopItems, AccountEmail: accountEmail, DeliveryAddress: deliveryAddress}
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
	es.BaseCommand
	ShopItems []*models.ShopItem `json:"shopItems" bson:"shopItems,omitempty" validate:"required"`
}

func NewOrderUpdatedCommand(aggregateID string, shopItems []*models.ShopItem) *OrderUpdatedCommand {
	return &OrderUpdatedCommand{BaseCommand: es.NewBaseCommand(aggregateID), ShopItems: shopItems}
}

type OrderCanceledCommand struct {
	es.BaseCommand
	CancelReason string `json:"cancelReason" validate:"required"`
}

func NewOrderCanceledCommand(aggregateID string, cancelReason string) *OrderCanceledCommand {
	return &OrderCanceledCommand{BaseCommand: es.NewBaseCommand(aggregateID), CancelReason: cancelReason}
}

type OrderDeliveredCommand struct {
	es.BaseCommand
	DeliveryTimestamp time.Time `json:"deliveryTimestamp" validate:"required"`
}

func NewOrderDeliveredCommand(aggregateID string, deliveryTimestamp time.Time) *OrderDeliveredCommand {
	return &OrderDeliveredCommand{BaseCommand: es.NewBaseCommand(aggregateID), DeliveryTimestamp: deliveryTimestamp}
}

type OrderChangeDeliveryAddressCommand struct {
	es.BaseCommand
	DeliveryAddress string `json:"deliveryAddress" bson:"deliveryAddress,omitempty" validate:"required"`
}

func NewOrderChangeDeliveryAddressCommand(aggregateID string, deliveryAddress string) *OrderChangeDeliveryAddressCommand {
	return &OrderChangeDeliveryAddressCommand{BaseCommand: es.NewBaseCommand(aggregateID), DeliveryAddress: deliveryAddress}
}
