package server

import (
	"context"
	"github.com/AleksK1NG/es-microservice/config"
	"github.com/AleksK1NG/es-microservice/internal/order/delivery/http"
	"github.com/AleksK1NG/es-microservice/internal/order/projection/elastic_projection"
	"github.com/AleksK1NG/es-microservice/internal/order/projection/mongo_projection"
	"github.com/AleksK1NG/es-microservice/internal/order/repository"
	"github.com/AleksK1NG/es-microservice/internal/order/service"
	"github.com/AleksK1NG/es-microservice/pkg/elasticsearch"
	"github.com/AleksK1NG/es-microservice/pkg/es/store"
	"github.com/AleksK1NG/es-microservice/pkg/eventstroredb"
	"github.com/AleksK1NG/es-microservice/pkg/interceptors"
	"github.com/AleksK1NG/es-microservice/pkg/logger"
	"github.com/AleksK1NG/es-microservice/pkg/middlewares"
	"github.com/AleksK1NG/es-microservice/pkg/mongodb"
	"github.com/AleksK1NG/es-microservice/pkg/tracing"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
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
	im            interceptors.InterceptorManager
	mw            middlewares.MiddlewareManager
	os            *service.OrderService
	v             *validator.Validate
	mongoClient   *mongo.Client
	elasticClient *v7.Client
	echo          *echo.Echo
}

func NewServer(cfg *config.Config, log logger.Logger) *server {
	return &server{cfg: cfg, log: log, v: validator.New(), echo: echo.New()}
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
	s.mw = middlewares.NewMiddlewareManager(s.log, s.cfg)

	mongoDBConn, err := mongodb.NewMongoDBConn(ctx, s.cfg.Mongo)
	if err != nil {
		return errors.Wrap(err, "NewMongoDBConn")
	}
	s.mongoClient = mongoDBConn
	defer mongoDBConn.Disconnect(ctx) // nolint: errcheck
	s.log.Infof("(Mongo connected) SessionsInProgress: {%v}", mongoDBConn.NumberSessionsInProgress())

	elasticClient, err := elasticsearch.NewElasticClient(s.cfg.Elastic)
	if err != nil {
		return err
	}
	s.elasticClient = elasticClient

	mongoRepository := repository.NewMongoRepository(s.log, s.cfg, s.mongoClient)
	elasticRepository := repository.NewElasticRepository(s.log, s.cfg, s.elasticClient)

	db, err := eventstroredb.NewEventStoreDB(s.cfg.EventStoreConfig)
	if err != nil {
		return err
	}

	aggregateStore := store.NewAggregateStore(s.log, db)
	s.os = service.NewOrderService(s.log, s.cfg, aggregateStore, mongoRepository, elasticRepository)

	mongoProjection := mongo_projection.NewOrderProjection(s.log, db, mongoRepository, s.cfg)
	elasticProjection := elastic_projection.NewElasticProjection(s.log, db, elasticRepository, s.cfg)

	go func() {
		err := mongoProjection.Subscribe(ctx, []string{s.cfg.Subscriptions.OrderPrefix}, s.cfg.Subscriptions.PoolSize, mongoProjection.ProcessEvents)
		if err != nil {
			s.log.Errorf("(orderProjection.Subscribe) err: {%v}", err)
			cancel()
		}
	}()

	go func() {
		err := elasticProjection.Subscribe(ctx, []string{s.cfg.Subscriptions.OrderPrefix}, s.cfg.Subscriptions.PoolSize, elasticProjection.ProcessEvents)
		if err != nil {
			s.log.Errorf("(elasticProjection.Subscribe) err: {%v}", err)
			cancel()
		}
	}()

	orderHandlers := http.NewOrderHandlers(s.echo.Group(s.cfg.Http.OrdersPath), s.log, s.mw, s.cfg, s.v, s.os)
	orderHandlers.MapRoutes()

	go func() {
		if err := s.runHttpServer(); err != nil {
			s.log.Errorf("(s.runHttpServer) err: {%v}", err)
			cancel()
		}
	}()
	s.log.Infof("(EventSourcingService) is listening on PORT: {%s}", s.cfg.Http.Port)

	closeGrpcServer, grpcServer, err := s.newOrderGrpcServer()
	if err != nil {
		cancel()
		return err
	}
	defer closeGrpcServer() // nolint: errcheck

	<-ctx.Done()
	grpcServer.GracefulStop()
	s.log.Info("server exited properly")
	return nil
}
