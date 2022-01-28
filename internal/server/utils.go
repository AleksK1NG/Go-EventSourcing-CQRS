package server

import (
	"context"
	"fmt"
	"github.com/AleksK1NG/es-microservice/config"
	"github.com/AleksK1NG/es-microservice/pkg/constants"
	"github.com/AleksK1NG/es-microservice/pkg/elasticsearch"
	serviceErrors "github.com/AleksK1NG/es-microservice/pkg/service_errors"
	"github.com/AleksK1NG/es-microservice/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

const (
	waitShotDownDuration = 3 * time.Second
)

func (s *server) initMongoDBCollections(ctx context.Context) {
	err := s.mongoClient.Database(s.cfg.Mongo.Db).CreateCollection(ctx, s.cfg.MongoCollections.Orders)
	if err != nil {
		if !utils.CheckErrMessages(err, serviceErrors.ErrMsgMongoCollectionAlreadyExists) {
			s.log.Warnf("(CreateCollection) err: {%v}", err)
		}
	}

	indexOptions := options.Index().SetSparse(true).SetUnique(true)
	index, err := s.mongoClient.Database(s.cfg.Mongo.Db).Collection(s.cfg.MongoCollections.Orders).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: constants.OrderIdIndex, Value: 1}},
		Options: indexOptions,
	})
	if err != nil && !utils.CheckErrMessages(err, serviceErrors.ErrMsgAlreadyExists) {
		s.log.Warnf("(CreateOne) err: {%v}", err)
	}
	s.log.Infof("(CreatedIndex) index: {%s}", index)

	list, err := s.mongoClient.Database(s.cfg.Mongo.Db).Collection(s.cfg.MongoCollections.Orders).Indexes().List(ctx)
	if err != nil {
		s.log.Warnf("(initMongoDBCollections) [List] err: {%v}", err)
	}

	if list != nil {
		var results []bson.M
		if err := list.All(ctx, &results); err != nil {
			s.log.Warnf("(All) err: {%v}", err)
		}
		s.log.Infof("(indexes) results: {%#v}", results)
	}

	collections, err := s.mongoClient.Database(s.cfg.Mongo.Db).ListCollectionNames(ctx, bson.M{})
	if err != nil {
		s.log.Warnf("(ListCollections) err: {%v}", err)
	}
	s.log.Infof("(Collections) created collections: {%v}", collections)
}

func (s *server) initElasticClient(ctx context.Context) error {
	elasticClient, err := elasticsearch.NewElasticClient(s.cfg.Elastic)
	if err != nil {
		return err
	}
	s.elasticClient = elasticClient

	info, code, err := s.elasticClient.Ping(s.cfg.Elastic.URL).Do(ctx)
	if err != nil {
		return errors.Wrap(err, "client.Ping")
	}
	s.log.Infof("Elasticsearch returned with code {%d} and version {%s}", code, info.Version.Number)

	esVersion, err := s.elasticClient.ElasticsearchVersion(s.cfg.Elastic.URL)
	if err != nil {
		return errors.Wrap(err, "client.ElasticsearchVersion")
	}
	s.log.Infof("Elasticsearch version {%s}", esVersion)

	return nil
}

func (s *server) runMetrics(cancel context.CancelFunc) {
	metricsServer := echo.New()
	go func() {
		metricsServer.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
			StackSize:         stackSize,
			DisablePrintStack: true,
			DisableStackAll:   true,
		}))
		metricsServer.GET(s.cfg.Probes.PrometheusPath, echo.WrapHandler(promhttp.Handler()))
		s.log.Infof("Metrics server is running on port: {%s}", s.cfg.Probes.PrometheusPort)
		if err := metricsServer.Start(s.cfg.Probes.PrometheusPort); err != nil {
			s.log.Errorf("metricsServer.Start: {%v}", err)
			cancel()
		}
	}()
}

func (s *server) getHttpMetricsCb() func(err error) {
	return func(err error) {
		if err != nil {
			s.metrics.ErrorHttpRequests.Inc()
		} else {
			s.metrics.SuccessHttpRequests.Inc()
		}
	}
}

func (s *server) getGrpcMetricsCb() func(err error) {
	return func(err error) {
		if err != nil {
			s.metrics.ErrorGrpcRequests.Inc()
		} else {
			s.metrics.SuccessGrpcRequests.Inc()
		}
	}
}

func (s *server) waitShootDown(duration time.Duration) {
	go func() {
		time.Sleep(duration)
		s.doneCh <- struct{}{}
	}()
}

func GetMicroserviceName(cfg *config.Config) string {
	return fmt.Sprintf("(%s)", strings.ToUpper(cfg.ServiceName))
}
