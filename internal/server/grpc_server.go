package server

import (
	grpc2 "github.com/AleksK1NG/es-microservice/internal/order/delivery/grpc"
	"github.com/AleksK1NG/es-microservice/pkg/constants"
	orderService "github.com/AleksK1NG/es-microservice/proto/order"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
)

const (
	maxConnectionIdle = 5
	gRPCTimeout       = 15
	maxConnectionAge  = 5
	gRPCTime          = 10
)

func (s *server) newOrderGrpcServer() (func() error, *grpc.Server, error) {
	l, err := net.Listen(constants.Tcp, s.cfg.GRPC.Port)
	if err != nil {
		return nil, nil, errors.Wrap(err, "net.Listen")
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: maxConnectionIdle * time.Minute,
			Timeout:           gRPCTimeout * time.Second,
			MaxConnectionAge:  maxConnectionAge * time.Minute,
			Time:              gRPCTime * time.Minute,
		}),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_recovery.UnaryServerInterceptor(),
			s.im.Logger,
		),
		),
	)

	grpcService := grpc2.NewOrderGrpcService(s.log, s.os, s.v, s.metrics)
	orderService.RegisterOrderServiceServer(grpcServer, grpcService)
	grpc_prometheus.Register(grpcServer)

	if s.cfg.GRPC.Development {
		reflection.Register(grpcServer)
	}

	go func() {
		s.log.Infof("%s gRPC server is listening on port: {%s}", GetMicroserviceName(s.cfg), s.cfg.GRPC.Port)
		s.log.Error(grpcServer.Serve(l))
	}()

	return l.Close, grpcServer, nil
}
