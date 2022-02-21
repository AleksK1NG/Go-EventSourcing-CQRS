package elastic_projection

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

func (o *elasticProjection) onOrderCreate(ctx context.Context, evt es.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "elasticProjection.onOrderCreate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData v1.OrderCreatedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	op := &models.OrderProjection{
		OrderID:      aggregate.GetOrderAggregateID(evt.AggregateID),
		ShopItems:    eventData.ShopItems,
		AccountEmail: eventData.AccountEmail,
		TotalPrice:   aggregate.GetShopItemsTotalPrice(eventData.ShopItems),
	}

	return o.elasticRepository.IndexOrder(ctx, op)
}

func (o *elasticProjection) onOrderPaid(ctx context.Context, evt es.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "elasticProjection.onOrderPaid")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var payment models.Payment
	if err := evt.GetJsonData(&payment); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	projection, err := o.elasticRepository.GetByID(ctx, aggregate.GetOrderAggregateID(evt.AggregateID))
	if err != nil {
		return err
	}
	projection.Paid = true
	projection.Payment = payment

	return o.elasticRepository.UpdateOrder(ctx, projection)
}

func (o *elasticProjection) onSubmit(ctx context.Context, evt es.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "elasticProjection.onSubmit")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	projection, err := o.elasticRepository.GetByID(ctx, aggregate.GetOrderAggregateID(evt.AggregateID))
	if err != nil {
		return err
	}
	projection.Submitted = true

	return o.elasticRepository.UpdateOrder(ctx, projection)
}

func (o *elasticProjection) onShoppingCartUpdate(ctx context.Context, evt es.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "elasticProjection.onShoppingCartUpdate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData v1.ShoppingCartUpdatedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	projection, err := o.elasticRepository.GetByID(ctx, aggregate.GetOrderAggregateID(evt.AggregateID))
	if err != nil {
		return err
	}
	projection.ShopItems = eventData.ShopItems
	projection.TotalPrice = aggregate.GetShopItemsTotalPrice(eventData.ShopItems)

	return o.elasticRepository.UpdateOrder(ctx, projection)
}

func (o *elasticProjection) onCancel(ctx context.Context, evt es.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "elasticProjection.onCancel")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData v1.OrderCanceledEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	projection, err := o.elasticRepository.GetByID(ctx, aggregate.GetOrderAggregateID(evt.AggregateID))
	if err != nil {
		return err
	}
	projection.Canceled = true
	projection.Completed = false
	projection.CancelReason = eventData.CancelReason

	return o.elasticRepository.UpdateOrder(ctx, projection)
}

func (o *elasticProjection) onComplete(ctx context.Context, evt es.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "elasticProjection.onComplete")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData v1.OrderCompletedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	projection, err := o.elasticRepository.GetByID(ctx, aggregate.GetOrderAggregateID(evt.AggregateID))
	if err != nil {
		return err
	}
	projection.Completed = true
	projection.DeliveredTime = eventData.DeliveryTimestamp

	return o.elasticRepository.UpdateOrder(ctx, projection)
}

func (o *elasticProjection) onDeliveryAddressChnaged(ctx context.Context, evt es.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "elasticProjection.onDeliveryAddressChnaged")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData v1.OrderDeliveryAddressChangedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	projection, err := o.elasticRepository.GetByID(ctx, aggregate.GetOrderAggregateID(evt.AggregateID))
	if err != nil {
		return err
	}
	projection.DeliveryAddress = eventData.DeliveryAddress

	return o.elasticRepository.UpdateOrder(ctx, projection)

}
