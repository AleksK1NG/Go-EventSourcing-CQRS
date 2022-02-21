package v1

type OrderCommands struct {
	CreateOrder                CreateOrderCommandHandler
	OrderPaid                  PayOrderCommandHandler
	SubmitOrder                SubmitOrderCommandHandler
	UpdateOrder                UpdateShoppingCartCommandHandler
	CancelOrder                CancelOrderCommandHandler
	CompleteOrder              CompleteOrderCommandHandler
	ChangeOrderDeliveryAddress ChangeDeliveryAddressCommandHandler
}

func NewOrderCommands(
	createOrder CreateOrderCommandHandler,
	orderPaid PayOrderCommandHandler,
	submitOrder SubmitOrderCommandHandler,
	updateOrder UpdateShoppingCartCommandHandler,
	cancelOrder CancelOrderCommandHandler,
	deliveryOrder CompleteOrderCommandHandler,
	changeOrderDeliveryAddress ChangeDeliveryAddressCommandHandler,
) *OrderCommands {
	return &OrderCommands{
		CreateOrder:                createOrder,
		OrderPaid:                  orderPaid,
		SubmitOrder:                submitOrder,
		UpdateOrder:                updateOrder,
		CancelOrder:                cancelOrder,
		CompleteOrder:              deliveryOrder,
		ChangeOrderDeliveryAddress: changeOrderDeliveryAddress,
	}
}
