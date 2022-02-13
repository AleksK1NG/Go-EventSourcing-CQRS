package queries

import (
	"context"
	"github.com/AleksK1NG/es-microservice/config"
	"github.com/AleksK1NG/es-microservice/internal/dto"
	"github.com/AleksK1NG/es-microservice/internal/order/repository"
	"github.com/AleksK1NG/es-microservice/pkg/es"
	"github.com/AleksK1NG/es-microservice/pkg/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type SearchOrdersQueryHandler interface {
	Handle(ctx context.Context, command *SearchOrdersQuery) (*dto.OrderSearchResponseDto, error)
}

type searchOrdersHandler struct {
	log               logger.Logger
	cfg               *config.Config
	es                es.AggregateStore
	elasticRepository repository.ElasticOrderRepository
}

func NewSearchOrdersHandler(log logger.Logger, cfg *config.Config, es es.AggregateStore, elasticRepository repository.ElasticOrderRepository) *searchOrdersHandler {
	return &searchOrdersHandler{log: log, cfg: cfg, es: es, elasticRepository: elasticRepository}
}

func (s *searchOrdersHandler) Handle(ctx context.Context, command *SearchOrdersQuery) (*dto.OrderSearchResponseDto, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "searchOrdersHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("SearchText", command.SearchText))

	return s.elasticRepository.Search(ctx, command.SearchText, command.Pq)
}
