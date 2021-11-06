package aggregate

import (
	"encoding/json"
	"github.com/AleksK1NG/es-microservice/internal/order/events"
	"github.com/pkg/errors"
)

func (a *OrderAggregate) handleCreateOrderCommand(command *CreateOrderCommand) error {
	createdData := &events.OrderCreatedData{ShopItems: command.ShopItems, AccountEmail: command.AccountEmail}
	createdDataBytes, err := json.Marshal(createdData)
	if err != nil {
		return err
	}

	createOrderEvent := events.NewCreateOrderEvent(a, createdDataBytes)

	if a.Order.Created || a.Version != 0 || command.OrderCreatedData.ShopItems == nil {
		return errors.New("already created")
	}

	return a.Apply(createOrderEvent)
}

func (a *OrderAggregate) handleOrderPaidCommand(command *OrderPaidCommand) error {

	if !a.Order.Created || a.Order.Canceled {
		return errors.New("order created or cancelled")
	}
	if a.Order.Paid {
		return errors.New("already paid")
	}
	if a.Order.Submitted {
		return errors.New("already submitted")
	}

	payOrderEvent := events.NewPayOrderEvent(a, nil)

	return a.Apply(payOrderEvent)
}

func (a *OrderAggregate) handleSubmitOrderCommand(command *SubmitOrderCommand) error {

	if !a.Order.Created || a.Order.Canceled {
		return errors.New("order created or cancelled")
	}
	if !a.Order.Paid {
		return errors.New("order not paid")
	}
	if a.Order.Submitted {
		return errors.New("already submitted")
	}

	submitOrderEvent := events.NewSubmitOrderEvent(a)

	return a.Apply(submitOrderEvent)
}

func (a *OrderAggregate) handleOrderUpdatedCommand(command *OrderUpdatedCommand) error {

	if !a.Order.Created || a.Order.Canceled {
		return errors.New("order created or cancelled")
	}
	if a.Order.Submitted {
		return errors.New("already submitted")
	}

	eventData := &events.OrderUpdatedData{ShopItems: command.ShopItems}
	eventDataBytes, err := json.Marshal(eventData)
	if err != nil {
		return err
	}

	orderUpdatedEvent := events.NewOrderUpdatedEvent(a, eventDataBytes)

	return a.Apply(orderUpdatedEvent)
}
