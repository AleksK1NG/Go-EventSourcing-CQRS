package dto

import (
	"time"
)

type OrderResponseDto struct {
	ID              string     `json:"id" bson:"_id,omitempty"`
	OrderID         string     `json:"orderId,omitempty" bson:"orderId,omitempty"`
	ShopItems       []ShopItem `json:"shopItems,omitempty" bson:"shopItems,omitempty"`
	AccountEmail    string     `json:"accountEmail,omitempty" bson:"accountEmail,omitempty" validate:"required,email"`
	DeliveryAddress string     `json:"deliveryAddress,omitempty" bson:"deliveryAddress,omitempty"`
	CancelReason    string     `json:"cancelReason,omitempty" bson:"cancelReason,omitempty"`
	TotalPrice      float64    `json:"totalPrice,omitempty" bson:"totalPrice,omitempty"`
	DeliveredTime   time.Time  `json:"deliveredTime,omitempty" bson:"deliveredTime,omitempty"`
	Created         bool       `json:"created,omitempty" bson:"created,omitempty"`
	Paid            bool       `json:"paid,omitempty" bson:"paid,omitempty"`
	Submitted       bool       `json:"submitted,omitempty" bson:"submitted,omitempty"`
	Completed       bool       `json:"completed,omitempty" bson:"completed,omitempty"`
	Canceled        bool       `json:"canceled,omitempty" bson:"canceled,omitempty"`
	Payment         Payment    `json:"payment,omitempty" bson:"payment,omitempty"`
}
