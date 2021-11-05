package models

import orderService "github.com/AleksK1NG/es-microservice/proto/order"

type Order struct {
	ItemsIDs   []string `json:"itemsIDs" bson:"itemsIDs,omitempty"`
	Created    bool     `json:"created" bson:"created,omitempty"`
	Paid       bool     `json:"paid" bson:"paid,omitempty"`
	Submitted  bool     `json:"submitted" bson:"submitted,omitempty"`
	Delivering bool     `json:"delivering" bson:"delivering,omitempty"`
	Delivered  bool     `json:"delivered" bson:"delivered,omitempty"`
	Canceled   bool     `json:"canceled" bson:"canceled,omitempty"`
}

func NewOrder() *Order {
	return &Order{
		ItemsIDs:   make([]string, 0),
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
		ItemsIDs:   order.ItemsIDs,
		Created:    order.Created,
		Paid:       order.Paid,
		Submitted:  order.Submitted,
		Delivering: order.Delivering,
		Delivered:  order.Delivered,
		Canceled:   order.Canceled,
	}
}

type OrderProjection struct {
	ID         string   `json:"id" bson:"_id,omitempty"`
	OrderID    string   `json:"orderId" bson:"orderId,omitempty"`
	ItemsIDs   []string `json:"itemsIDs" bson:"itemsIDs,omitempty"`
	Created    bool     `json:"created" bson:"created,omitempty"`
	Paid       bool     `json:"paid" bson:"paid,omitempty"`
	Submitted  bool     `json:"submitted" bson:"submitted,omitempty"`
	Delivering bool     `json:"delivering" bson:"delivering,omitempty"`
	Delivered  bool     `json:"delivered" bson:"delivered,omitempty"`
	Canceled   bool     `json:"canceled" bson:"canceled,omitempty"`
}

func OrderProjectionToProto(order *OrderProjection) *orderService.Order {
	return &orderService.Order{
		ID:         order.OrderID,
		ItemsIDs:   order.ItemsIDs,
		Created:    order.Created,
		Paid:       order.Paid,
		Submitted:  order.Submitted,
		Delivering: order.Delivering,
		Delivered:  order.Delivered,
		Canceled:   order.Canceled,
	}
}
