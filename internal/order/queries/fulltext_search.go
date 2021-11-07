package queries

import (
	"context"
	"github.com/AleksK1NG/es-microservice/config"
	"github.com/AleksK1NG/es-microservice/internal/order/repository"
	"github.com/AleksK1NG/es-microservice/pkg/es"
	"github.com/AleksK1NG/es-microservice/pkg/logger"
	orderService "github.com/AleksK1NG/es-microservice/proto/order"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type SearchOrdersQueryHandler interface {
	Handle(ctx context.Context, command *SearchOrdersQuery) (*orderService.SearchRes, error)
}

type searchOrdersHandler struct {
	log               logger.Logger
	cfg               *config.Config
	es                es.AggregateStore
	elasticRepository repository.ElasticRepository
}

func NewSearchOrdersHandler(log logger.Logger, cfg *config.Config, es es.AggregateStore, elasticRepository repository.ElasticRepository) *searchOrdersHandler {
	return &searchOrdersHandler{log: log, cfg: cfg, es: es, elasticRepository: elasticRepository}
}

func (s *searchOrdersHandler) Handle(ctx context.Context, command *SearchOrdersQuery) (*orderService.SearchRes, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "searchOrdersHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("SearchText", command.SearchText))

	return s.elasticRepository.Search(ctx, command.SearchText, command.Pq)
}
