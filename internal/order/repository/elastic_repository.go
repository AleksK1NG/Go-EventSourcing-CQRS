package repository

import (
	"context"
	"encoding/json"
	"github.com/AleksK1NG/es-microservice/config"
	"github.com/AleksK1NG/es-microservice/internal/models"
	"github.com/AleksK1NG/es-microservice/pkg/logger"
	"github.com/AleksK1NG/es-microservice/pkg/tracing"
	v7 "github.com/olivere/elastic/v7"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type ElasticRepository interface {
	IndexOrder(ctx context.Context, order *models.OrderProjection) error
	GetByID(ctx context.Context, orderID string) (*models.OrderProjection, error)
	UpdateOrder(ctx context.Context, order *models.OrderProjection) error
	Search(ctx context.Context, text string) ([]*models.OrderProjection, error)
}

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
		return err
	}

	e.log.Infof("IndexOrder result: %+v", res)
	return nil
}

func (e *elasticRepository) GetByID(ctx context.Context, orderID string) (*models.OrderProjection, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "elasticRepository.GetByID")
	defer span.Finish()
	span.LogFields(log.String("OrderID", orderID))

	result, err := e.elasticClient.Get().Index(e.cfg.ElasticIndexes.Orders).Id(orderID).FetchSource(true).Do(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	jsonData, err := result.Source.MarshalJSON()
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	var order models.OrderProjection
	if err := json.Unmarshal(jsonData, &order); err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &order, nil
}

func (e *elasticRepository) UpdateOrder(ctx context.Context, order *models.OrderProjection) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "elasticRepository.UpdateOrder")
	defer span.Finish()
	span.LogFields(log.String("OrderID", order.OrderID))

	res, err := e.elasticClient.Update().Index(e.cfg.ElasticIndexes.Orders).Id(order.OrderID).Doc(order).FetchSource(true).Do(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	e.log.Infof("UpdateOrder result: %+v", res)
	return nil
}

func (e *elasticRepository) Search(ctx context.Context, text string) ([]*models.OrderProjection, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "elasticRepository.Search")
	defer span.Finish()
	span.LogFields(log.String("Search", text))

	shouldMatch := v7.NewBoolQuery().
		Should(
			v7.NewMatchPhrasePrefixQuery("shopItems.title", text),
			v7.NewMatchPhrasePrefixQuery("shopItems.description", text),
		).
		MinimumNumberShouldMatch(1)

	searchResult, err := e.elasticClient.Search(e.cfg.ElasticIndexes.Orders).
		Query(shouldMatch).
		From(0).
		Explain(true).
		FetchSource(true).
		Version(true).
		//Sort("price", true).
		Size(60).     // take documents 0-9
		Pretty(true). // pretty print request and response JSON
		Do(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	orders := make([]*models.OrderProjection, 0, len(searchResult.Hits.Hits))
	for _, hit := range searchResult.Hits.Hits {
		jsonBytes, err := hit.Source.MarshalJSON()
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		var order models.OrderProjection
		if err := json.Unmarshal(jsonBytes, &order); err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		orders = append(orders, &order)
	}

	return orders, nil
}
