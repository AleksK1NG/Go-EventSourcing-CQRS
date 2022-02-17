package v1

import (
	"time"

	"github.com/AleksK1NG/es-microservice/internal/order/models"
)

type OrderCreatedEvent struct {
	ShopItems       []*models.ShopItem `json:"shopItems" bson:"shopItems,omitempty" validate:"required"`
	AccountEmail    string             `json:"accountEmail" bson:"accountEmail,omitempty" validate:"required,email"`
	DeliveryAddress string             `json:"deliveryAddress" bson:"deliveryAddress,omitempty" validate:"required"`
}

type OrderUpdatedEvent struct {
	ShopItems []*models.ShopItem `json:"shopItems" bson:"shopItems,omitempty" validate:"required"`
}

type OrderCanceledEvent struct {
	CancelReason string `json:"cancelReason" validate:"required"`
}

type OrderDeliveredEvent struct {
	DeliveryTimestamp time.Time `json:"deliveryTimestamp" validate:"required"`
}

type OrderDeliveryAddressChangedEvent struct {
	DeliveryAddress string `json:"deliveryAddress" bson:"deliveryAddress,omitempty" validate:"required"`
}
