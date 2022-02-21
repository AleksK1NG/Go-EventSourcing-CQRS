package dto

import "github.com/AleksK1NG/es-microservice/internal/order/models"

type UpdateShoppingItemsReqDto struct {
	ShopItems []*models.ShopItem `json:"shopItems" bson:"shopItems,omitempty" validate:"required"`
}
