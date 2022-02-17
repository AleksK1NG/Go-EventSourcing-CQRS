package aggregate

import (
	"context"
	"time"

	eventsV1 "github.com/AleksK1NG/es-microservice/internal/order/events/v1"
	"github.com/AleksK1NG/es-microservice/internal/order/models"
	"github.com/AleksK1NG/es-microservice/pkg/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

func (a *OrderAggregate) CreateOrder(ctx context.Context, shopItems []*models.ShopItem, accountEmail, deliveryAddress string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrderAggregate.CreateOrder")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", a.GetID()))

	if a.Order.Created {
		return ErrAlreadyCreated
	}
	if shopItems == nil {
		return ErrOrderShopItemsIsRequired
	}
	if deliveryAddress == "" {
		return ErrInvalidDeliveryAddress
	}

	event, err := eventsV1.NewOrderCreatedEvent(a, shopItems, accountEmail, deliveryAddress)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrderCreatedEvent")
	}

	if err := event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "SetMetadata")
	}

	return a.Apply(event)
}

func (a *OrderAggregate) PayOrder(ctx context.Context, payment models.Payment) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrderAggregate.PayOrder")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", a.GetID()))

	if !a.Order.Created || a.Order.Canceled {
		return ErrAlreadyCreatedOrCancelled
	}
	if a.Order.Paid {
		return ErrAlreadyPaid
	}
	if a.Order.Submitted {
		return ErrAlreadySubmitted
	}

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

func (a *OrderAggregate) SubmitOrder(ctx context.Context) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrderAggregate.SubmitOrder")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", a.GetID()))

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

func (a *OrderAggregate) UpdateOrder(ctx context.Context, shopItems []*models.ShopItem) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrderAggregate.UpdateOrder")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", a.GetID()))

	if !a.Order.Created || a.Order.Canceled {
		return ErrAlreadyCreatedOrCancelled
	}
	if a.Order.Submitted {
		return ErrAlreadySubmitted
	}

	orderUpdatedEvent, err := eventsV1.NewOrderUpdatedEvent(a, shopItems)
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

func (a *OrderAggregate) CancelOrder(ctx context.Context, cancelReason string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrderAggregate.CancelOrder")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", a.GetID()))

	if a.Order.Delivered {
		return ErrOrderAlreadyDelivered
	}
	if cancelReason == "" {
		return ErrCancelReasonRequired
	}

	event, err := eventsV1.NewOrderCanceledEvent(a, cancelReason)
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

func (a *OrderAggregate) DeliverOrder(ctx context.Context, deliveryTimestamp time.Time) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrderAggregate.DeliverOrder")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", a.GetID()))

	if a.Order.Delivered {
		return ErrOrderAlreadyDelivered
	}
	if a.Order.Canceled {
		return ErrOrderAlreadyCanceled
	}
	if !a.Order.Paid {
		return ErrOrderMustBePaidBeforeDelivered
	}

	event, err := eventsV1.NewOrderDeliveredEvent(a, deliveryTimestamp)
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

func (a *OrderAggregate) ChangeDeliveryAddress(ctx context.Context, deliveryAddress string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrderAggregate.ChangeDeliveryAddress")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", a.GetID()))

	if a.Order.Delivered {
		return ErrOrderAlreadyDelivered
	}

	event, err := eventsV1.NewOrderDeliveryAddressChangedEvent(a, deliveryAddress)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrderDeliveryAddressChangedEvent")
	}

	if err := event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "SetMetadata")
	}

	return a.Apply(event)
}
