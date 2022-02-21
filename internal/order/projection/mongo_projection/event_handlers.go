package mongo_projection

import (
	"context"

	"github.com/AleksK1NG/es-microservice/internal/order/aggregate"
	"github.com/AleksK1NG/es-microservice/internal/order/events/v1"
	"github.com/AleksK1NG/es-microservice/internal/order/models"
	"github.com/AleksK1NG/es-microservice/pkg/es"
	"github.com/AleksK1NG/es-microservice/pkg/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

func (o *mongoProjection) onOrderCreate(ctx context.Context, evt es.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoProjection.onOrderCreate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData v1.OrderCreatedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	span.LogFields(log.String("AccountEmail", eventData.AccountEmail))

	op := &models.OrderProjection{
		OrderID:         aggregate.GetOrderAggregateID(evt.AggregateID),
		ShopItems:       eventData.ShopItems,
		AccountEmail:    eventData.AccountEmail,
		TotalPrice:      aggregate.GetShopItemsTotalPrice(eventData.ShopItems),
		DeliveryAddress: eventData.DeliveryAddress,
	}

	_, err := o.mongoRepo.Insert(ctx, op)
	if err != nil {
		return err
	}

	return nil
}

func (o *mongoProjection) onOrderPaid(ctx context.Context, evt es.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoProjection.onOrderPaid")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var payment models.Payment
	if err := evt.GetJsonData(&payment); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	op := &models.OrderProjection{OrderID: aggregate.GetOrderAggregateID(evt.AggregateID), Paid: true, Payment: payment}
	return o.mongoRepo.UpdatePayment(ctx, op)
}

func (o *mongoProjection) onSubmit(ctx context.Context, evt es.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoProjection.onSubmit")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	op := &models.OrderProjection{OrderID: aggregate.GetOrderAggregateID(evt.AggregateID), Submitted: true}
	return o.mongoRepo.UpdateSubmit(ctx, op)
}

func (o *mongoProjection) onShoppingCartUpdate(ctx context.Context, evt es.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoProjection.onShoppingCartUpdate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData v1.ShoppingCartUpdatedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	op := &models.OrderProjection{OrderID: aggregate.GetOrderAggregateID(evt.AggregateID), ShopItems: eventData.ShopItems}
	op.TotalPrice = aggregate.GetShopItemsTotalPrice(eventData.ShopItems)
	return o.mongoRepo.UpdateOrder(ctx, op)
}

func (o *mongoProjection) onCancel(ctx context.Context, evt es.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoProjection.onCancel")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData v1.OrderCanceledEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	op := &models.OrderProjection{
		OrderID:      aggregate.GetOrderAggregateID(evt.AggregateID),
		Canceled:     true,
		Completed:    false,
		CancelReason: eventData.CancelReason,
	}
	return o.mongoRepo.UpdateCancel(ctx, op)
}

func (o *mongoProjection) onCompleted(ctx context.Context, evt es.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoProjection.onCompleted")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData v1.OrderCompletedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	op := &models.OrderProjection{
		OrderID:       aggregate.GetOrderAggregateID(evt.AggregateID),
		Canceled:      false,
		Completed:     true,
		DeliveredTime: eventData.DeliveryTimestamp,
	}
	return o.mongoRepo.Complete(ctx, op)
}

func (o *mongoProjection) onDeliveryAddressChnaged(ctx context.Context, evt es.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoProjection.onDeliveryAddressChnaged")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData v1.OrderDeliveryAddressChangedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	op := &models.OrderProjection{
		OrderID:         aggregate.GetOrderAggregateID(evt.AggregateID),
		DeliveryAddress: eventData.DeliveryAddress,
	}
	return o.mongoRepo.UpdateDeliveryAddress(ctx, op)
}
