package aggregate

import (
	"context"

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
	if command.ShopItems == nil {
		return ErrOrderShopItemsIsRequired
	}
	if command.DeliveryAddress == "" {
		return ErrInvalidDeliveryAddress
	}

	event, err := eventsV1.NewCreateOrderEvent(a, command.ShopItems, command.AccountEmail, command.DeliveryAddress)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewCreateOrderEvent")
	}

	if err := event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "SetMetadata")
	}

	return a.Apply(event)
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
	event, err := eventsV1.NewOrderPaidEvent(a, &payment)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrderPaidEvent")
	}

	if err := event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "SetMetadata")
	}

	return a.Apply(event)
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

	submitOrderEvent, err := eventsV1.NewSubmitOrderEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewSubmitOrderEvent")
	}

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

	orderUpdatedEvent, err := eventsV1.NewOrderUpdatedEvent(a, command.ShopItems)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrderUpdatedEvent")
	}

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

	event, err := eventsV1.NewOrderCanceledEvent(a, command.CancelReason)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrderCanceledEvent")
	}

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

	event, err := eventsV1.NewOrderDeliveredEvent(a, command.DeliveryTimestamp)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrderDeliveredEvent")
	}

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

	event, err := eventsV1.NewOrderDeliveryAddressUpdatedEvent(a, command.DeliveryAddress)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrderDeliveryAddressUpdatedEvent")
	}

	if err := event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "SetMetadata")
	}

	return a.Apply(event)
}
