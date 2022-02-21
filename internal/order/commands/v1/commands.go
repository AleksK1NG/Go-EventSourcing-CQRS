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

type PayOrderCommand struct {
	models.Payment
	es.BaseCommand
}

func NewPayOrderCommand(payment models.Payment, aggregateID string) *PayOrderCommand {
	return &PayOrderCommand{Payment: payment, BaseCommand: es.NewBaseCommand(aggregateID)}
}

type SubmitOrderCommand struct {
	es.BaseCommand
}

func NewSubmitOrderCommand(aggregateID string) *SubmitOrderCommand {
	return &SubmitOrderCommand{BaseCommand: es.NewBaseCommand(aggregateID)}
}

type UpdateShoppingCartCommand struct {
	es.BaseCommand
	ShopItems []*models.ShopItem `json:"shopItems" bson:"shopItems,omitempty" validate:"required"`
}

func NewUpdateShoppingCartCommand(aggregateID string, shopItems []*models.ShopItem) *UpdateShoppingCartCommand {
	return &UpdateShoppingCartCommand{BaseCommand: es.NewBaseCommand(aggregateID), ShopItems: shopItems}
}

type CancelOrderCommand struct {
	es.BaseCommand
	CancelReason string `json:"cancelReason" validate:"required"`
}

func NewCancelOrderCommand(aggregateID string, cancelReason string) *CancelOrderCommand {
	return &CancelOrderCommand{BaseCommand: es.NewBaseCommand(aggregateID), CancelReason: cancelReason}
}

type CompleteOrderCommand struct {
	es.BaseCommand
	DeliveryTimestamp time.Time `json:"deliveryTimestamp" validate:"required"`
}

func NewCompleteOrderCommand(aggregateID string, deliveryTimestamp time.Time) *CompleteOrderCommand {
	return &CompleteOrderCommand{BaseCommand: es.NewBaseCommand(aggregateID), DeliveryTimestamp: deliveryTimestamp}
}

type ChangeDeliveryAddressCommand struct {
	es.BaseCommand
	DeliveryAddress string `json:"deliveryAddress" bson:"deliveryAddress,omitempty" validate:"required"`
}

func NewChangeDeliveryAddressCommand(aggregateID string, deliveryAddress string) *ChangeDeliveryAddressCommand {
	return &ChangeDeliveryAddressCommand{BaseCommand: es.NewBaseCommand(aggregateID), DeliveryAddress: deliveryAddress}
}
