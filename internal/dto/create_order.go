package dto

import "github.com/AleksK1NG/es-microservice/internal/order/models"

type CreateOrderReqDto struct {
	ShopItems       []*models.ShopItem `json:"shopItems" bson:"shopItems,omitempty" validate:"required"`
	AccountEmail    string             `json:"accountEmail" bson:"accountEmail,omitempty" validate:"required,email"`
	DeliveryAddress string             `json:"deliveryAddress" bson:"deliveryAddress,omitempty" validate:"required"`
}
