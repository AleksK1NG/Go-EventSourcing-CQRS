package events

import (
	"github.com/AleksK1NG/es-microservice/internal/models"
	"time"
)

type OrderCreatedEventData struct {
	ShopItems    []*models.ShopItem `json:"shopItems" bson:"shopItems,omitempty" validate:"required"`
	AccountEmail string             `json:"accountEmail" bson:"accountEmail,omitempty" validate:"required,email"`
}

type OrderUpdatedEventData struct {
	ShopItems []*models.ShopItem `json:"shopItems" bson:"shopItems,omitempty" validate:"required"`
}

type OrderCanceledEventData struct {
	CancelReason string `json:"cancelReason"`
}

type OrderDeliveredEventData struct {
	DeliveryTimestamp time.Time `json:"deliveryTimestamp"`
}
