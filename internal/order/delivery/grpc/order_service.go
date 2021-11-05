package grpc

import (
	"context"
	"github.com/AleksK1NG/es-microservice/internal/models"
	"github.com/AleksK1NG/es-microservice/internal/order/aggregate"
	"github.com/AleksK1NG/es-microservice/internal/order/events"
	"github.com/AleksK1NG/es-microservice/internal/order/queries"
	"github.com/AleksK1NG/es-microservice/internal/order/service"
	"github.com/AleksK1NG/es-microservice/pkg/logger"
	"github.com/AleksK1NG/es-microservice/proto/order"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type orderGrpcService struct {
	log logger.Logger
	os  *service.OrderService
}

func NewOrderGrpcService(log logger.Logger, os *service.OrderService) *orderGrpcService {
	return &orderGrpcService{log: log, os: os}
}

func (s *orderGrpcService) CreateOrder(ctx context.Context, req *orderService.CreateOrderReq) (*orderService.CreateOrderRes, error) {
	command := aggregate.NewCreateOrderCommand(events.OrderCreatedData{ItemsIDs: req.GetItemID()}, req.GetAggregateID())
	if err := s.os.Commands.CreateOrder.Handle(ctx, command); err != nil {
		s.log.WarnMsg("CreateOrder.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}

	return &orderService.CreateOrderRes{AggregateID: req.GetAggregateID()}, nil
}

func (s *orderGrpcService) PayOrder(ctx context.Context, req *orderService.PayOrderReq) (*orderService.PayOrderRes, error) {
	command := aggregate.NewOrderPaidCommand(req.GetAggregateID())
	if err := s.os.Commands.OrderPaid.Handle(ctx, command); err != nil {
		s.log.WarnMsg("OrderPaid.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}

	return &orderService.PayOrderRes{AggregateID: req.GetAggregateID()}, nil
}

func (s *orderGrpcService) SubmitOrder(ctx context.Context, req *orderService.SubmitOrderReq) (*orderService.SubmitOrderRes, error) {
	command := aggregate.NewSubmitOrderCommand(req.GetAggregateID())
	if err := s.os.Commands.SubmitOrder.Handle(ctx, command); err != nil {
		s.log.WarnMsg("SubmitOrder.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}

	return &orderService.SubmitOrderRes{AggregateID: req.GetAggregateID()}, nil
}

func (s *orderGrpcService) GetOrderByID(ctx context.Context, req *orderService.GetOrderByIDReq) (*orderService.GetOrderByIDRes, error) {
	query := queries.NewGetOrderByIDQuery(req.GetAggregateID())
	orderProjection, err := s.os.Queries.GetOrderByIDQuery.Handle(ctx, query)
	if err != nil {
		s.log.WarnMsg("GetOrderByIDQuery.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}

	return &orderService.GetOrderByIDRes{Order: models.OrderProjectionToProto(orderProjection)}, nil
}

func (s *orderGrpcService) UpdateOrder(ctx context.Context, req *orderService.UpdateOrderReq) (*orderService.UpdateOrderRes, error) {
	command := aggregate.NewOrderUpdatedCommand(events.OrderCreatedData{ItemsIDs: req.GetItemID()}, req.GetAggregateID())
	if err := s.os.Commands.UpdateOrder.Handle(ctx, command); err != nil {
		return nil, err
	}

	return &orderService.UpdateOrderRes{}, nil
}

func (s *orderGrpcService) errResponse(c codes.Code, err error) error {
	return status.Error(c, err.Error())
}
