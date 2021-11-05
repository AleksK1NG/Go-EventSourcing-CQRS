package commands

type OrderCommands struct {
	CreateOrder CreateOrderCommandHandler
	OrderPaid   OrderPaidCommandHandler
	SubmitOrder SubmitOrderCommandHandler
	UpdateOrder UpdateOrderCommandHandler
}

func NewOrderCommands(
	createOrder CreateOrderCommandHandler,
	orderPaid OrderPaidCommandHandler,
	submitOrder SubmitOrderCommandHandler,
	updateOrder UpdateOrderCommandHandler,
) *OrderCommands {
	return &OrderCommands{CreateOrder: createOrder, OrderPaid: orderPaid, SubmitOrder: submitOrder, UpdateOrder: updateOrder}
}
