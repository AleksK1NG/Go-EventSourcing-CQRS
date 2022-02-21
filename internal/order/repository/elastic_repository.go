package repository

import (
	"context"
	"encoding/json"

	"github.com/AleksK1NG/es-microservice/config"
	"github.com/AleksK1NG/es-microservice/internal/dto"
	"github.com/AleksK1NG/es-microservice/internal/mappers"
	"github.com/AleksK1NG/es-microservice/internal/order/models"
	"github.com/AleksK1NG/es-microservice/pkg/logger"
	"github.com/AleksK1NG/es-microservice/pkg/tracing"
	"github.com/AleksK1NG/es-microservice/pkg/utils"
	v7 "github.com/olivere/elastic/v7"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

const (
	shopItemTitle            = "shopItems.title"
	shopItemDescription      = "shopItems.description"
	minimumNumberShouldMatch = 1
)

type elasticRepository struct {
	log           logger.Logger
	cfg           *config.Config
	elasticClient *v7.Client
}

func NewElasticRepository(log logger.Logger, cfg *config.Config, elasticClient *v7.Client) *elasticRepository {
	return &elasticRepository{log: log, cfg: cfg, elasticClient: elasticClient}
}

func (e *elasticRepository) IndexOrder(ctx context.Context, order *models.OrderProjection) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "elasticRepository.IndexOrder")
	defer span.Finish()
	span.LogFields(log.String("OrderID", order.OrderID))

	res, err := e.elasticClient.Index().Index(e.cfg.ElasticIndexes.Orders).BodyJson(order).Id(order.OrderID).Do(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "elasticClient.Index")
	}

	e.log.Debugf("(IndexOrder) result: {%s}", res.Result)
	return nil
}

func (e *elasticRepository) GetByID(ctx context.Context, orderID string) (*models.OrderProjection, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "elasticRepository.GetByID")
	defer span.Finish()
	span.LogFields(log.String("OrderID", orderID))

	result, err := e.elasticClient.Get().Index(e.cfg.ElasticIndexes.Orders).Id(orderID).FetchSource(true).Do(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "elasticClient.Get")
	}

	jsonData, err := result.Source.MarshalJSON()
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "Source.MarshalJSON")
	}

	var order models.OrderProjection
	if err := json.Unmarshal(jsonData, &order); err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "json.Unmarshal")
	}

	return &order, nil
}

func (e *elasticRepository) UpdateOrder(ctx context.Context, order *models.OrderProjection) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "elasticRepository.UpdateShoppingCart")
	defer span.Finish()
	span.LogFields(log.String("OrderID", order.OrderID))

	res, err := e.elasticClient.Update().Index(e.cfg.ElasticIndexes.Orders).Id(order.OrderID).Doc(order).FetchSource(true).Do(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "elasticClient.Update")
	}

	e.log.Debugf("(UpdateShoppingCart) result: {%s}", res.Result)
	return nil
}

func (e *elasticRepository) Search(ctx context.Context, text string, pq *utils.Pagination) (*dto.OrderSearchResponseDto, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "elasticRepository.Search")
	defer span.Finish()
	span.LogFields(log.String("Search", text))

	shouldMatch := v7.NewBoolQuery().
		Should(v7.NewMatchPhrasePrefixQuery(shopItemTitle, text), v7.NewMatchPhrasePrefixQuery(shopItemDescription, text)).
		MinimumNumberShouldMatch(minimumNumberShouldMatch)

	searchResult, err := e.elasticClient.Search(e.cfg.ElasticIndexes.Orders).
		Query(shouldMatch).
		From(pq.GetOffset()).
		Explain(e.cfg.Elastic.Explain).
		FetchSource(e.cfg.Elastic.FetchSource).
		Version(e.cfg.Elastic.Version).
		Size(pq.GetSize()).
		Pretty(e.cfg.Elastic.Pretty).
		Do(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "elasticClient.Search")
	}

	orders := make([]*models.OrderProjection, 0, len(searchResult.Hits.Hits))
	for _, hit := range searchResult.Hits.Hits {
		jsonBytes, err := hit.Source.MarshalJSON()
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, errors.Wrap(err, "Source.MarshalJSON")
		}
		var order models.OrderProjection
		if err := json.Unmarshal(jsonBytes, &order); err != nil {
			tracing.TraceErr(span, err)
			return nil, errors.Wrap(err, "json.Unmarshal")
		}
		orders = append(orders, &order)
	}

	return &dto.OrderSearchResponseDto{
		Pagination: dto.Pagination{
			TotalCount: searchResult.TotalHits(),
			TotalPages: int64(pq.GetTotalPages(int(searchResult.TotalHits()))),
			Page:       int64(pq.GetPage()),
			Size:       int64(pq.GetSize()),
			HasMore:    pq.GetHasMore(int(searchResult.TotalHits())),
		},
		Orders: mappers.OrdersFromProjections(orders),
	}, nil
}
