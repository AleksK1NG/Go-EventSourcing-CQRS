package aggregate

import (
	"context"
	"encoding/json"
	"github.com/AleksK1NG/es-microservice/internal/dto"
	"github.com/AleksK1NG/es-microservice/internal/order/events"
	serviceErrors "github.com/AleksK1NG/es-microservice/pkg/service_errors"
	"github.com/AleksK1NG/es-microservice/pkg/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func (a *OrderAggregate) onCreateOrderCommand(ctx context.Context, command *CreateOrderCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrderAggregate.onCreateOrderCommand")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", command.GetAggregateID()))

	if a.Order.Created || a.Version != 0 || command.OrderCreatedData.ShopItems == nil {
		return serviceErrors.ErrAlreadyCreatedOrCancelled
	}

	createdData := &dto.OrderCreatedData{ShopItems: command.ShopItems, AccountEmail: command.AccountEmail}
	createdDataBytes, err := json.Marshal(createdData)
	if err != nil {
		return err
	}

	createOrderEvent := events.NewCreateOrderEvent(a, createdDataBytes)
	if err := createOrderEvent.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		return err
	}

	return a.Apply(createOrderEvent)
}

func (a *OrderAggregate) onOrderPaidCommand(ctx context.Context, command *OrderPaidCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrderAggregate.onOrderPaidCommand")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", command.GetAggregateID()))

	if !a.Order.Created || a.Order.Canceled {
		return serviceErrors.ErrAlreadyCreatedOrCancelled
	}
	if a.Order.Paid {
		return serviceErrors.ErrAlreadyPaid
	}
	if a.Order.Submitted {
		return serviceErrors.ErrAlreadySubmitted
	}

	payOrderEvent := events.NewPayOrderEvent(a, nil)
	if err := payOrderEvent.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		return err
	}

	return a.Apply(payOrderEvent)
}

func (a *OrderAggregate) onSubmitOrderCommand(ctx context.Context, command *SubmitOrderCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrderAggregate.onSubmitOrderCommand")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", command.GetAggregateID()))

	if !a.Order.Created || a.Order.Canceled {
		return serviceErrors.ErrAlreadyCreatedOrCancelled
	}
	if !a.Order.Paid {
		return serviceErrors.ErrOrderNotPaid
	}
	if a.Order.Submitted {
		return serviceErrors.ErrAlreadySubmitted
	}

	submitOrderEvent := events.NewSubmitOrderEvent(a)
	if err := submitOrderEvent.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		return err
	}

	return a.Apply(submitOrderEvent)
}

func (a *OrderAggregate) onOrderUpdatedCommand(ctx context.Context, command *OrderUpdatedCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrderAggregate.onOrderUpdatedCommand")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", command.GetAggregateID()))

	if !a.Order.Created || a.Order.Canceled {
		return serviceErrors.ErrAlreadyCreatedOrCancelled
	}
	if a.Order.Submitted {
		return serviceErrors.ErrAlreadySubmitted
	}

	eventData := &dto.OrderUpdatedData{ShopItems: command.ShopItems}
	eventDataBytes, err := json.Marshal(eventData)
	if err != nil {
		return err
	}

	orderUpdatedEvent := events.NewOrderUpdatedEvent(a, eventDataBytes)
	if err := orderUpdatedEvent.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		return err
	}

	return a.Apply(orderUpdatedEvent)
}
