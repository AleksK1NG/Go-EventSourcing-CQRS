package aggregate

import (
	"context"
	"encoding/json"
	"github.com/AleksK1NG/es-microservice/internal/order/commands/v1"
	eventsV1 "github.com/AleksK1NG/es-microservice/internal/order/events/v1"
	"github.com/AleksK1NG/es-microservice/internal/order/models"
	"github.com/AleksK1NG/es-microservice/pkg/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

func (a *OrderAggregate) CreateOrder(ctx context.Context, command *v1.CreateOrderCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrderAggregate.CreateOrder")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", command.GetAggregateID()))

	if a.Order.Created {
		return ErrAlreadyCreated
	}
	if command.OrderCreatedEventData.ShopItems == nil {
		return ErrOrderShopItemsIsRequired
	}
	if command.DeliveryAddress == "" {
		return ErrInvalidDeliveryAddress
	}

	createdData := &eventsV1.OrderCreatedEventData{ShopItems: command.ShopItems, AccountEmail: command.AccountEmail, DeliveryAddress: command.DeliveryAddress}
	createdDataBytes, err := json.Marshal(createdData)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "json.Marshal")
	}

	createOrderEvent := eventsV1.NewCreateOrderEvent(a, createdDataBytes)
	if err := createOrderEvent.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "SetMetadata")
	}

	return a.Apply(createOrderEvent)
}

func (a *OrderAggregate) PayOrder(ctx context.Context, command *v1.OrderPaidCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrderAggregate.PayOrder")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", command.GetAggregateID()))

	if !a.Order.Created || a.Order.Canceled {
		return ErrAlreadyCreatedOrCancelled
	}
	if a.Order.Paid {
		return ErrAlreadyPaid
	}
	if a.Order.Submitted {
		return ErrAlreadySubmitted
	}

	payment := models.Payment{PaymentID: command.PaymentID, Timestamp: command.Timestamp}
	eventData, err := json.Marshal(&payment)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "json.Marshal")
	}

	payOrderEvent := eventsV1.NewPayOrderEvent(a, eventData)
	if err := payOrderEvent.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "SetMetadata")
	}

	return a.Apply(payOrderEvent)
}

func (a *OrderAggregate) SubmitOrder(ctx context.Context, command *v1.SubmitOrderCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrderAggregate.SubmitOrder")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", command.GetAggregateID()))

	if !a.Order.Created || a.Order.Canceled {
		return ErrAlreadyCreatedOrCancelled
	}
	if !a.Order.Paid {
		return ErrOrderNotPaid
	}
	if a.Order.Submitted {
		return ErrAlreadySubmitted
	}

	submitOrderEvent := eventsV1.NewSubmitOrderEvent(a)
	if err := submitOrderEvent.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "SetMetadata")
	}

	return a.Apply(submitOrderEvent)
}

func (a *OrderAggregate) UpdateOrder(ctx context.Context, command *v1.OrderUpdatedCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrderAggregate.UpdateOrder")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", command.GetAggregateID()))

	if !a.Order.Created || a.Order.Canceled {
		return ErrAlreadyCreatedOrCancelled
	}
	if a.Order.Submitted {
		return ErrAlreadySubmitted
	}

	eventData := &eventsV1.OrderUpdatedEventData{ShopItems: command.ShopItems}
	eventDataBytes, err := json.Marshal(eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "json.Marshal")
	}

	orderUpdatedEvent := eventsV1.NewOrderUpdatedEvent(a, eventDataBytes)
	if err := orderUpdatedEvent.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "SetMetadata")
	}

	return a.Apply(orderUpdatedEvent)
}

func (a *OrderAggregate) CancelOrder(ctx context.Context, command *v1.OrderCanceledCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrderAggregate.CancelOrder")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", command.GetAggregateID()))

	if a.Order.Delivered {
		return ErrOrderAlreadyDelivered
	}
	if command.CancelReason == "" {
		return ErrCancelReasonRequired
	}

	eventData := &eventsV1.OrderCanceledEventData{CancelReason: command.CancelReason}
	eventDataBytes, err := json.Marshal(eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "json.Marshal")
	}

	event := eventsV1.NewOrderCanceledEvent(a, eventDataBytes)
	if err := event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "SetMetadata")
	}

	return a.Apply(event)
}

func (a *OrderAggregate) DeliverOrder(ctx context.Context, command *v1.OrderDeliveredCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrderAggregate.DeliverOrder")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", command.GetAggregateID()))

	if a.Order.Delivered {
		return ErrOrderAlreadyDelivered
	}
	if a.Order.Canceled {
		return ErrOrderAlreadyCanceled
	}
	if !a.Order.Paid {
		return ErrOrderMustBePaidBeforeDelivered
	}

	eventData := &eventsV1.OrderDeliveredEventData{DeliveryTimestamp: command.DeliveryTimestamp}
	eventDataBytes, err := json.Marshal(eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "json.Marshal")
	}

	event := eventsV1.NewOrderDeliveredEvent(a, eventDataBytes)
	if err := event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "SetMetadata")
	}

	return a.Apply(event)
}

func (a *OrderAggregate) ChangeDeliveryAddress(ctx context.Context, command *v1.OrderChangeDeliveryAddressCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrderAggregate.ChangeDeliveryAddress")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", command.GetAggregateID()))

	if a.Order.Delivered {
		return ErrOrderAlreadyDelivered
	}

	eventData := &eventsV1.OrderChangeDeliveryAddress{DeliveryAddress: command.DeliveryAddress}
	eventDataBytes, err := json.Marshal(eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "json.Marshal")
	}

	event := eventsV1.NewOrderDeliveryAddressUpdatedEvent(a, eventDataBytes)
	if err := event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "SetMetadata")
	}

	return a.Apply(event)
}
