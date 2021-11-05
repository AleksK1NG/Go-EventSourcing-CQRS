package grpc

import (
	"context"
	"github.com/AleksK1NG/es-microservice/pkg/logger"
	"github.com/AleksK1NG/es-microservice/proto/order"
)

type orderGrpcService struct {
	log logger.Logger
}

func NewOrderGrpcService(log logger.Logger) *orderGrpcService {
	return &orderGrpcService{log: log}
}

func (o *orderGrpcService) CreateOrder(ctx context.Context, req *orderService.CreateOrderReq) (*orderService.CreateOrderRes, error) {
	//TODO implement me
	panic("implement me")
}

func (o *orderGrpcService) PayOrder(ctx context.Context, req *orderService.PayOrderReq) (*orderService.PayOrderRes, error) {
	//TODO implement me
	panic("implement me")
}

func (o *orderGrpcService) SubmitOrder(ctx context.Context, req *orderService.SubmitOrderReq) (*orderService.SubmitOrderRes, error) {
	//TODO implement me
	panic("implement me")
}

func (o *orderGrpcService) GetOrderByID(ctx context.Context, req *orderService.GetOrderByIDReq) (*orderService.GetOrderByIDRes, error) {
	//TODO implement me
	panic("implement me")
}

func (o *orderGrpcService) UpdateOrder(ctx context.Context, req *orderService.UpdateOrderReq) (*orderService.UpdateOrderRes, error) {
	//TODO implement me
	panic("implement me")
}
