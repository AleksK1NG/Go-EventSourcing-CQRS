package dto

import "github.com/AleksK1NG/es-microservice/internal/models"

type OrderCreatedData struct {
	ShopItems    []*models.ShopItem `json:"shopItems" bson:"shopItems,omitempty" validate:"required"`
	AccountEmail string             `json:"accountEmail" bson:"accountEmail,omitempty" validate:"required,email"`
}

type OrderUpdatedData struct {
	ShopItems []*models.ShopItem `json:"shopItems" bson:"shopItems,omitempty" validate:"required"`
}
