package mongo_projection

import (
	"context"
	"github.com/AleksK1NG/es-microservice/internal/models"
	"github.com/AleksK1NG/es-microservice/internal/order/aggregate"
	"github.com/AleksK1NG/es-microservice/internal/order/events"
	"github.com/AleksK1NG/es-microservice/pkg/es"
	"github.com/AleksK1NG/es-microservice/pkg/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func (o *orderProjection) handleOrderCreateEvent(ctx context.Context, evt es.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "orderProjection.handleOrderCreateEvent")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrderCreatedData
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	op := &models.OrderProjection{
		OrderID:      aggregate.GetOrderAggregateID(evt.AggregateID),
		ShopItems:    eventData.ShopItems,
		Created:      true,
		Paid:         false,
		Submitted:    false,
		Delivering:   false,
		Delivered:    false,
		Canceled:     false,
		AccountEmail: eventData.AccountEmail,
		TotalPrice:   aggregate.GetShopItemsTotalPrice(eventData.ShopItems),
	}

	result, err := o.mongoRepo.Insert(ctx, op)
	if err != nil {
		return err
	}

	o.log.Debugf("projection OrderCreated result: %s", result)
	return nil
}

func (o *orderProjection) handleOrderPaidEvent(ctx context.Context, evt es.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "orderProjection.handleOrderPaidEvent")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	op := &models.OrderProjection{OrderID: aggregate.GetOrderAggregateID(evt.AggregateID), Paid: true}
	return o.mongoRepo.UpdateOrder(ctx, op)
}

func (o *orderProjection) handleSubmitEvent(ctx context.Context, evt es.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "orderProjection.handleSubmitEvent")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	op := &models.OrderProjection{OrderID: aggregate.GetOrderAggregateID(evt.AggregateID), Submitted: true}
	return o.mongoRepo.UpdateOrder(ctx, op)
}

func (o *orderProjection) handleUpdateEvent(ctx context.Context, evt es.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "orderProjection.handleUpdateEvent")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrderUpdatedData
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	op := &models.OrderProjection{OrderID: aggregate.GetOrderAggregateID(evt.AggregateID), ShopItems: eventData.ShopItems}
	op.TotalPrice = aggregate.GetShopItemsTotalPrice(eventData.ShopItems)
	return o.mongoRepo.UpdateOrder(ctx, op)
}
