package models

import (
	"fmt"
	orderService "github.com/AleksK1NG/es-microservice/proto/order"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Order struct {
	ID              string      `json:"id" bson:"_id,omitempty"`
	ShopItems       []*ShopItem `json:"shopItems" bson:"shopItems,omitempty"`
	AccountEmail    string      `json:"accountEmail" bson:"accountEmail,omitempty"`
	DeliveryAddress string      `json:"deliveryAddress" bson:"deliveryAddress,omitempty"`
	CancelReason    string      `json:"cancelReason" bson:"cancelReason,omitempty"`
	TotalPrice      float64     `json:"totalPrice" bson:"totalPrice,omitempty"`
	DeliveredTime   time.Time   `json:"deliveredTime" bson:"deliveredTime,omitempty"`
	Created         bool        `json:"created" bson:"created,omitempty"`
	Paid            bool        `json:"paid" bson:"paid,omitempty"`
	Submitted       bool        `json:"submitted" bson:"submitted,omitempty"`
	Delivered       bool        `json:"delivered" bson:"delivered,omitempty"`
	Canceled        bool        `json:"canceled" bson:"canceled,omitempty"`
	Payment         Payment     `json:"payment" bson:"payment,omitempty"`
}

func (o *Order) String() string {
	return fmt.Sprintf("ID: {%s}, ShopItems: {%+v}, Created: {%v}, Paid: {%v}, Submitted: {%v}, "+
		"Delivered: {%v}, Canceled: {%v}, CancelReason: {%s}, TotalPrice: {%v}, AccountEmail: {%s}, DeliveryAddress: {%s}, DeliveredTime: {%s}, Payment: {%s}",
		o.ID,
		o.ShopItems,
		o.Created,
		o.Paid,
		o.Submitted,
		o.Delivered,
		o.Canceled,
		o.CancelReason,
		o.TotalPrice,
		o.AccountEmail,
		o.DeliveryAddress,
		o.DeliveredTime.UTC().String(),
		o.Payment.String(),
	)
}

func NewOrder() *Order {
	return &Order{
		ShopItems: make([]*ShopItem, 0),
		Created:   false,
		Paid:      false,
		Submitted: false,
		Delivered: false,
		Canceled:  false,
	}
}

func OrderToProto(order *Order, id string) *orderService.Order {
	return &orderService.Order{
		ID:                id,
		ShopItems:         ShopItemsToProto(order.ShopItems),
		Created:           order.Created,
		Paid:              order.Paid,
		Submitted:         order.Submitted,
		Delivered:         order.Delivered,
		Canceled:          order.Canceled,
		CancelReason:      order.CancelReason,
		DeliveryTimestamp: timestamppb.New(order.DeliveredTime),
		DeliveryAddress:   order.DeliveryAddress,
		AccountEmail:      order.AccountEmail,
		TotalPrice:        order.TotalPrice,
		Payment:           PaymentToProto(order.Payment),
	}
}
