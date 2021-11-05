package queries

import (
	"context"
	"github.com/AleksK1NG/es-microservice/config"
	"github.com/AleksK1NG/es-microservice/internal/mappers"
	"github.com/AleksK1NG/es-microservice/internal/models"
	"github.com/AleksK1NG/es-microservice/internal/order/aggregate"
	"github.com/AleksK1NG/es-microservice/pkg/es"
	"github.com/AleksK1NG/es-microservice/pkg/logger"
	"github.com/pkg/errors"
)

type GetOrderByIDQueryHandler interface {
	Handle(ctx context.Context, command *GetOrderByIDQuery) (*models.OrderProjection, error)
}

type getOrderByIDHandler struct {
	log logger.Logger
	cfg *config.Config
	es  es.AggregateStore
	//mongoRepo repository.OrderRepository
}

func NewGetOrderByIDHandler(log logger.Logger, cfg *config.Config, es es.AggregateStore) *getOrderByIDHandler {
	return &getOrderByIDHandler{log: log, cfg: cfg, es: es}
}

func (q *getOrderByIDHandler) Handle(ctx context.Context, command *GetOrderByIDQuery) (*models.OrderProjection, error) {
	//orderProjection, err := q.mongoRepo.GetByID(ctx, command.ID)
	//if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
	//	return nil, err
	//}
	//if orderProjection != nil {
	//	return orderProjection, nil
	//}

	order := aggregate.NewOrderAggregateWithID(command.ID)
	if err := q.es.Load(ctx, order); err != nil {
		return nil, err
	}

	if len(order.AppliedEvents) == 0 || order.GetVersion() == 0 {
		return nil, errors.New("order not found")
	}

	orderProjection := mappers.OrderProjectionFromAggregate(order)

	//_, err = q.mongoRepo.Insert(ctx, orderProjection)
	//if err != nil {
	//	return nil, err
	//}

	return orderProjection, nil
}
