package models

import (
	"fmt"
	orderService "github.com/AleksK1NG/es-microservice/proto/order"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type OrderProjection struct {
	ID              string      `json:"id" bson:"_id,omitempty"`
	OrderID         string      `json:"orderId,omitempty" bson:"orderId,omitempty"`
	ShopItems       []*ShopItem `json:"shopItems,omitempty" bson:"shopItems,omitempty"`
	AccountEmail    string      `json:"accountEmail,omitempty" bson:"accountEmail,omitempty" validate:"required,email"`
	DeliveryAddress string      `json:"deliveryAddress,omitempty" bson:"deliveryAddress,omitempty"`
	CancelReason    string      `json:"cancelReason,omitempty" bson:"cancelReason,omitempty"`
	TotalPrice      float64     `json:"totalPrice,omitempty" bson:"totalPrice,omitempty"`
	DeliveredTime   time.Time   `json:"deliveredTime,omitempty" bson:"deliveredTime,omitempty"`
	Created         bool        `json:"created,omitempty" bson:"created,omitempty"`
	Paid            bool        `json:"paid,omitempty" bson:"paid,omitempty"`
	Submitted       bool        `json:"submitted,omitempty" bson:"submitted,omitempty"`
	Delivered       bool        `json:"delivered,omitempty" bson:"delivered,omitempty"`
	Canceled        bool        `json:"canceled,omitempty" bson:"canceled,omitempty"`
	Payment         Payment     `json:"payment,omitempty" bson:"payment,omitempty"`
}

func (o *OrderProjection) String() string {
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

func OrderProjectionToProto(order *OrderProjection) *orderService.Order {
	return &orderService.Order{
		ID:                order.OrderID,
		ShopItems:         ShopItemsToProto(order.ShopItems),
		Created:           order.Created,
		Paid:              order.Paid,
		Submitted:         order.Submitted,
		Delivered:         order.Delivered,
		Canceled:          order.Canceled,
		TotalPrice:        order.TotalPrice,
		AccountEmail:      order.AccountEmail,
		CancelReason:      order.CancelReason,
		DeliveryTimestamp: timestamppb.New(order.DeliveredTime),
		DeliveryAddress:   order.DeliveryAddress,
		Payment:           PaymentToProto(order.Payment),
	}
}

func OrderProjectionsToProto(orderProjections []*OrderProjection) []*orderService.Order {
	orders := make([]*orderService.Order, 0, len(orderProjections))
	for _, projection := range orderProjections {
		orders = append(orders, OrderProjectionToProto(projection))
	}
	return orders
}

func PaymentToProto(payment Payment) *orderService.Payment {
	return &orderService.Payment{
		ID:        payment.PaymentID,
		Timestamp: timestamppb.New(payment.Timestamp),
	}
}

func PaymentFromProto(payment *orderService.Payment) Payment {
	return Payment{
		PaymentID: payment.GetID(),
		Timestamp: payment.GetTimestamp().AsTime(),
	}
}
