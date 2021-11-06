package server

import (
	"context"
	"github.com/AleksK1NG/es-microservice/config"
	"github.com/AleksK1NG/es-microservice/internal/order/projection"
	"github.com/AleksK1NG/es-microservice/internal/order/repository"
	"github.com/AleksK1NG/es-microservice/internal/order/service"
	"github.com/AleksK1NG/es-microservice/pkg/elasticsearch"
	"github.com/AleksK1NG/es-microservice/pkg/es/store"
	"github.com/AleksK1NG/es-microservice/pkg/eventstroredb"
	"github.com/AleksK1NG/es-microservice/pkg/interceptors"
	"github.com/AleksK1NG/es-microservice/pkg/logger"
	"github.com/AleksK1NG/es-microservice/pkg/mongodb"
	"github.com/AleksK1NG/es-microservice/pkg/tracing"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/go-playground/validator"
	v7 "github.com/olivere/elastic/v7"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"os/signal"
	"syscall"
)

type server struct {
	cfg           *config.Config
	log           logger.Logger
	db            *esdb.Client
	im            interceptors.InterceptorManager
	os            *service.OrderService
	v             *validator.Validate
	mongoClient   *mongo.Client
	elasticClient *v7.Client
}

func NewServer(cfg *config.Config, log logger.Logger) *server {
	return &server{cfg: cfg, log: log, v: validator.New()}
}

func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	if err := s.v.StructCtx(ctx, s.cfg); err != nil {
		return errors.Wrap(err, "cfg validate")
	}

	if s.cfg.Jaeger.Enable {
		tracer, closer, err := tracing.NewJaegerTracer(s.cfg.Jaeger)
		if err != nil {
			return err
		}
		defer closer.Close() // nolint: errcheck
		opentracing.SetGlobalTracer(tracer)
	}

	s.im = interceptors.NewInterceptorManager(s.log)

	mongoDBConn, err := mongodb.NewMongoDBConn(ctx, s.cfg.Mongo)
	if err != nil {
		return errors.Wrap(err, "NewMongoDBConn")
	}
	s.mongoClient = mongoDBConn
	defer mongoDBConn.Disconnect(ctx) // nolint: errcheck
	s.log.Infof("Mongo connected: %v", mongoDBConn.NumberSessionsInProgress())

	elasticClient, err := elasticsearch.NewElasticClient()
	if err != nil {
		return err
	}
	s.elasticClient = elasticClient

	mongoRepository := repository.NewMongoRepository(s.log, s.cfg, s.mongoClient)

	db, err := eventstroredb.NewEventStoreDB(s.cfg.EventStoreConfig)
	if err != nil {
		return err
	}

	aggregateStore := store.NewAggregateStore(s.log, db)
	s.os = service.NewOrderService(s.log, s.cfg, aggregateStore, mongoRepository)

	orderProjection := projection.NewOrderProjection(s.log, db, mongoRepository)

	go func() {
		err := orderProjection.Subscribe(ctx, []string{s.cfg.Subscriptions.OrderPrefix}, s.cfg.Subscriptions.PoolSize, orderProjection.ProcessEvents)
		if err != nil {
			s.log.Errorf("orderProjection.Subscribe: %v", err)
			cancel()
		}
	}()

	closeGrpcServer, grpcServer, err := s.newOrderGrpcServer()
	if err != nil {
		cancel()
		return err
	}
	defer closeGrpcServer() // nolint: errcheck

	<-ctx.Done()
	grpcServer.GracefulStop()
	s.log.Info("Order server exited properly")
	return nil
}
