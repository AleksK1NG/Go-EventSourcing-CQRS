package commands

type OrderCommands struct {
	CreateOrder   CreateOrderCommandHandler
	OrderPaid     OrderPaidCommandHandler
	SubmitOrder   SubmitOrderCommandHandler
	UpdateOrder   UpdateOrderCommandHandler
	CancelOrder   CancelOrderCommandHandler
	DeliveryOrder DeliveryOrderCommandHandler
}

func NewOrderCommands(
	createOrder CreateOrderCommandHandler,
	orderPaid OrderPaidCommandHandler,
	submitOrder SubmitOrderCommandHandler,
	updateOrder UpdateOrderCommandHandler,
	cancelOrder CancelOrderCommandHandler,
	deliveryOrder DeliveryOrderCommandHandler,
) *OrderCommands {
	return &OrderCommands{
		CreateOrder:   createOrder,
		OrderPaid:     orderPaid,
		SubmitOrder:   submitOrder,
		UpdateOrder:   updateOrder,
		CancelOrder:   cancelOrder,
		DeliveryOrder: deliveryOrder,
	}
}
