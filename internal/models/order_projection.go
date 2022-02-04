package models

import (
	"fmt"
	orderService "github.com/AleksK1NG/es-microservice/proto/order"
	"time"
)

type OrderProjection struct {
	ID              string      `json:"id" bson:"_id,omitempty"`
	OrderID         string      `json:"orderId" bson:"orderId,omitempty"`
	ShopItems       []*ShopItem `json:"shopItems" bson:"shopItems,omitempty"`
	AccountEmail    string      `json:"accountEmail" bson:"accountEmail,omitempty" validate:"required,email"`
	DeliveryAddress string      `json:"deliveryAddress" bson:"deliveryAddress,omitempty"`
	CancelReason    string      `json:"cancelReason" bson:"cancelReason,omitempty"`
	TotalPrice      float64     `json:"totalPrice" bson:"totalPrice,omitempty"`
	DeliveredTime   time.Time   `json:"deliveredTime" bson:"deliveredTime,omitempty"`
	Created         bool        `json:"created" bson:"created,omitempty"`
	Paid            bool        `json:"paid" bson:"paid,omitempty"`
	Submitted       bool        `json:"submitted" bson:"submitted,omitempty"`
	Delivering      bool        `json:"delivering" bson:"delivering,omitempty"`
	Delivered       bool        `json:"delivered" bson:"delivered,omitempty"`
	Canceled        bool        `json:"canceled" bson:"canceled,omitempty"`
}

func (o *OrderProjection) String() string {
	return fmt.Sprintf("ID: {%s}, ShopItems: {%+v}, Created: {%v}, Paid: {%v}, Submitted: {%v}, Delivering: {%v}, Delivered: {%v}, Canceled: {%v}, TotalPrice: {%v}, AccountEmail: {%s},",
		o.ID,
		o.ShopItems,
		o.Created,
		o.Paid,
		o.Submitted,
		o.Delivering,
		o.Delivered,
		o.Canceled,
		o.TotalPrice,
		o.AccountEmail,
	)
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

func OrderProjectionsToProto(orderProjections []*OrderProjection) []*orderService.Order {
	orders := make([]*orderService.Order, 0, len(orderProjections))
	for _, projection := range orderProjections {
		orders = append(orders, OrderProjectionToProto(projection))
	}
	return orders
}
