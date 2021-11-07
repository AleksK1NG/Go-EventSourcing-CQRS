package elastic_projection

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

func (o *elasticProjection) handleOrderCreateEvent(ctx context.Context, evt es.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "elasticProjection.handleOrderCreateEvent")
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

	return o.elasticRepository.IndexOrder(ctx, op)
}

func (o *elasticProjection) handleOrderPaidEvent(ctx context.Context, evt es.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "elasticProjection.handleOrderPaidEvent")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	op := &models.OrderProjection{OrderID: aggregate.GetOrderAggregateID(evt.AggregateID), Paid: true}
	return o.elasticRepository.UpdateOrder(ctx, op)
}

func (o *elasticProjection) handleSubmitEvent(ctx context.Context, evt es.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "elasticProjection.handleSubmitEvent")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	op := &models.OrderProjection{OrderID: aggregate.GetOrderAggregateID(evt.AggregateID), Submitted: true}
	return o.elasticRepository.UpdateOrder(ctx, op)
}

func (o *elasticProjection) handleUpdateEvent(ctx context.Context, evt es.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "elasticProjection.handleUpdateEvent")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrderUpdatedData
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	op := &models.OrderProjection{OrderID: aggregate.GetOrderAggregateID(evt.AggregateID), ShopItems: eventData.ShopItems}
	op.TotalPrice = aggregate.GetShopItemsTotalPrice(eventData.ShopItems)
	return o.elasticRepository.UpdateOrder(ctx, op)
}
