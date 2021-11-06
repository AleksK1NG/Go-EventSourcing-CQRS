package models

import (
	orderService "github.com/AleksK1NG/es-microservice/proto/order"
)

type Order struct {
	ID           string      `json:"id" bson:"_id,omitempty"`
	ShopItems    []*ShopItem `json:"shopItems" bson:"shopItems,omitempty"`
	Created      bool        `json:"created" bson:"created,omitempty"`
	Paid         bool        `json:"paid" bson:"paid,omitempty"`
	Submitted    bool        `json:"submitted" bson:"submitted,omitempty"`
	Delivering   bool        `json:"delivering" bson:"delivering,omitempty"`
	Delivered    bool        `json:"delivered" bson:"delivered,omitempty"`
	Canceled     bool        `json:"canceled" bson:"canceled,omitempty"`
	TotalPrice   float64     `json:"totalPrice" bson:"totalPrice,omitempty"`
	AccountEmail string      `json:"accountEmail" bson:"accountEmail,omitempty"`
}

func NewOrder() *Order {
	return &Order{
		ShopItems:  make([]*ShopItem, 0),
		Created:    false,
		Paid:       false,
		Submitted:  false,
		Delivering: false,
		Delivered:  false,
		Canceled:   false,
	}
}

func OrderToProto(order *Order, id string) *orderService.Order {
	return &orderService.Order{
		ID:         id,
		ShopItems:  ShopItemsToProto(order.ShopItems),
		Created:    order.Created,
		Paid:       order.Paid,
		Submitted:  order.Submitted,
		Delivering: order.Delivering,
		Delivered:  order.Delivered,
		Canceled:   order.Canceled,
	}
}

type OrderProjection struct {
	ID           string      `json:"id" bson:"_id,omitempty"`
	OrderID      string      `json:"orderId" bson:"orderId,omitempty"`
	ShopItems    []*ShopItem `json:"shopItems" bson:"shopItems,omitempty"`
	Created      bool        `json:"created" bson:"created,omitempty"`
	Paid         bool        `json:"paid" bson:"paid,omitempty"`
	Submitted    bool        `json:"submitted" bson:"submitted,omitempty"`
	Delivering   bool        `json:"delivering" bson:"delivering,omitempty"`
	Delivered    bool        `json:"delivered" bson:"delivered,omitempty"`
	Canceled     bool        `json:"canceled" bson:"canceled,omitempty"`
	TotalPrice   float64     `json:"totalPrice" bson:"totalPrice,omitempty"`
	AccountEmail string      `json:"accountEmail" bson:"accountEmail,omitempty" validate:"required,email"`
}

func OrderProjectionToProto(order *OrderProjection) *orderService.Order {
	return &orderService.Order{
		ID:           order.OrderID,
		ShopItems:    ShopItemsToProto(order.ShopItems),
		Created:      order.Created,
		Paid:         order.Paid,
		Submitted:    order.Submitted,
		Delivering:   order.Delivering,
		Delivered:    order.Delivered,
		Canceled:     order.Canceled,
		TotalPrice:   order.TotalPrice,
		AccountEmail: order.AccountEmail,
	}
}
