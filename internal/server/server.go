package server

import (
	"context"
	"github.com/AleksK1NG/es-microservice/config"
	"github.com/AleksK1NG/es-microservice/internal/order/projection"
	"github.com/AleksK1NG/es-microservice/internal/order/service"
	"github.com/AleksK1NG/es-microservice/pkg/es/store"
	"github.com/AleksK1NG/es-microservice/pkg/eventstroredb"
	"github.com/AleksK1NG/es-microservice/pkg/interceptors"
	"github.com/AleksK1NG/es-microservice/pkg/logger"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/go-playground/validator"
	"os"
	"os/signal"
	"syscall"
)

type server struct {
	cfg *config.Config
	log logger.Logger
	db  *esdb.Client
	im  interceptors.InterceptorManager
	os  *service.OrderService
	v   *validator.Validate
}

func NewServer(cfg *config.Config, log logger.Logger) *server {
	return &server{cfg: cfg, log: log, v: validator.New()}
}

func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	s.im = interceptors.NewInterceptorManager(s.log)

	//orderProjection := projection.NewOrderProjection(s.log, s.cfg, eventStore, mongoRepository)

	db, err := eventstroredb.NewEventStoreDB(s.cfg.EventStoreConfig)
	if err != nil {
		return err
	}

	aggregateStore := store.NewAggregateStore(s.log, db)
	s.os = service.NewOrderService(s.log, s.cfg, aggregateStore)

	orderProjection := projection.NewOrderProjection(s.log, db)

	go func() {
		s.log.Fatal(orderProjection.ProcessEvents(ctx))
	}()

	closeGrpcServer, grpcServer, err := s.newOrderGrpcServer()
	if err != nil {
		return err
	}
	defer closeGrpcServer() // nolint: errcheck

	<-ctx.Done()
	grpcServer.GracefulStop()
	s.log.Info("Order server exited properly")
	return nil
}
