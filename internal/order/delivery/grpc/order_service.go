package grpc

import (
	"context"
	"time"

	"github.com/AleksK1NG/es-microservice/internal/mappers"
	"github.com/AleksK1NG/es-microservice/internal/metrics"
	"github.com/AleksK1NG/es-microservice/internal/order/commands/v1"
	"github.com/AleksK1NG/es-microservice/internal/order/models"
	"github.com/AleksK1NG/es-microservice/internal/order/queries"
	"github.com/AleksK1NG/es-microservice/internal/order/service"
	grpcErrors "github.com/AleksK1NG/es-microservice/pkg/grpc_errors"
	"github.com/AleksK1NG/es-microservice/pkg/logger"
	"github.com/AleksK1NG/es-microservice/pkg/tracing"
	"github.com/AleksK1NG/es-microservice/pkg/utils"
	"github.com/AleksK1NG/es-microservice/proto/order"
	"github.com/go-playground/validator"
	"github.com/opentracing/opentracing-go/log"
	uuid "github.com/satori/go.uuid"
)

type orderGrpcService struct {
	log     logger.Logger
	os      *service.OrderService
	v       *validator.Validate
	metrics *metrics.ESMicroserviceMetrics
}

func NewOrderGrpcService(log logger.Logger, os *service.OrderService, v *validator.Validate, metrics *metrics.ESMicroserviceMetrics) *orderGrpcService {
	return &orderGrpcService{log: log, os: os, v: v, metrics: metrics}
}

func (s *orderGrpcService) CreateOrder(ctx context.Context, req *orderService.CreateOrderReq) (*orderService.CreateOrderRes, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "orderGrpcService.CreateOrder")
	defer span.Finish()
	span.LogFields(log.String("req", req.String()))
	s.metrics.CreateOrderGrpcRequests.Inc()

	aggregateID := uuid.NewV4().String()
	command := v1.NewCreateOrderCommand(aggregateID, models.ShopItemsFromProto(req.GetShopItems()), req.GetAccountEmail(), req.GetDeliveryAddress())
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.Errorf("(validate) aggregateID: {%s}, err: {%v}", aggregateID, err)
		tracing.TraceErr(span, err)
		return nil, s.errResponse(err)
	}

	if err := s.os.Commands.CreateOrder.Handle(ctx, command); err != nil {
		s.log.Errorf("(CreateOrder.Handle) orderID: {%s}, err: {%v}", aggregateID, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("(created order): orderID: {%s}", aggregateID)
	return &orderService.CreateOrderRes{AggregateID: aggregateID}, nil
}

func (s *orderGrpcService) PayOrder(ctx context.Context, req *orderService.PayOrderReq) (*orderService.PayOrderRes, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "orderGrpcService.PayOrder")
	defer span.Finish()
	span.LogFields(log.String("req", req.String()))
	s.metrics.PayOrderGrpcRequests.Inc()

	payment := models.Payment{PaymentID: req.GetPayment().GetID(), Timestamp: time.Now()}
	command := v1.NewPayOrderCommand(payment, req.GetAggregateID())
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.Errorf("(validate) err: {%v}", err)
		tracing.TraceErr(span, err)
		return nil, s.errResponse(err)
	}

	if err := s.os.Commands.OrderPaid.Handle(ctx, command); err != nil {
		s.log.Errorf("(OrderPaid.Handle) orderID: {%s}, err: {%v}", req.GetAggregateID(), err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("(paid order): orderID: {%s}", req.GetAggregateID())
	return &orderService.PayOrderRes{AggregateID: req.GetAggregateID()}, nil
}

func (s *orderGrpcService) SubmitOrder(ctx context.Context, req *orderService.SubmitOrderReq) (*orderService.SubmitOrderRes, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "orderGrpcService.SubmitOrder")
	defer span.Finish()
	span.LogFields(log.String("req", req.String()))
	s.metrics.SubmitOrderGrpcRequests.Inc()

	command := v1.NewSubmitOrderCommand(req.GetAggregateID())
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.Errorf("(validate) err: {%v}", err)
		tracing.TraceErr(span, err)
		return nil, s.errResponse(err)
	}

	if err := s.os.Commands.SubmitOrder.Handle(ctx, command); err != nil {
		s.log.Errorf("(SubmitOrder.Handle) orderID: {%s}, err: {%v}", req.GetAggregateID(), err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("(submitted order): orderID: {%s}", req.GetAggregateID())
	return &orderService.SubmitOrderRes{AggregateID: req.GetAggregateID()}, nil
}

func (s *orderGrpcService) GetOrderByID(ctx context.Context, req *orderService.GetOrderByIDReq) (*orderService.GetOrderByIDRes, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "orderGrpcService.GetOrderByID")
	defer span.Finish()
	span.LogFields(log.String("req", req.String()))
	s.metrics.GetOrderByIdGrpcRequests.Inc()

	query := queries.NewGetOrderByIDQuery(req.GetAggregateID())
	if err := s.v.StructCtx(ctx, query); err != nil {
		s.log.Errorf("(validate) err: {%v}", err)
		tracing.TraceErr(span, err)
		return nil, s.errResponse(err)
	}

	orderProjection, err := s.os.Queries.GetOrderByID.Handle(ctx, query)
	if err != nil {
		s.log.Errorf("(GetOrderByID.Handle) orderID: {%s}, err: {%v}", req.GetAggregateID(), err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("(GetOrderByID) AggregateID: {%s}", req.GetAggregateID())
	s.log.Debugf("(GetOrderByID) orderProjection: {%s}", orderProjection.String())
	return &orderService.GetOrderByIDRes{Order: models.OrderProjectionToProto(orderProjection)}, nil
}

func (s *orderGrpcService) UpdateShoppingCart(ctx context.Context, req *orderService.UpdateShoppingCartReq) (*orderService.UpdateShoppingCartRes, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "orderGrpcService.UpdateShoppingCart")
	defer span.Finish()
	span.LogFields(log.String("UpdateShoppingCart req", req.String()))
	s.metrics.UpdateOrderGrpcRequests.Inc()

	command := v1.NewUpdateShoppingCartCommand(req.GetAggregateID(), models.ShopItemsFromProto(req.GetShopItems()))
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.Errorf("(validate) err: {%v}", err)
		tracing.TraceErr(span, err)
		return nil, s.errResponse(err)
	}

	if err := s.os.Commands.UpdateOrder.Handle(ctx, command); err != nil {
		s.log.Errorf("(UpdateShoppingCart.Handle) orderID: {%s}, err: {%v}", req.GetAggregateID(), err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("(UpdateShoppingCart): AggregateID: {%s}", req.GetAggregateID())
	return &orderService.UpdateShoppingCartRes{}, nil
}

func (s *orderGrpcService) CancelOrder(ctx context.Context, req *orderService.CancelOrderReq) (*orderService.CancelOrderRes, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "orderGrpcService.CancelOrder")
	defer span.Finish()
	span.LogFields(log.String("CancelOrder req", req.String()))
	s.metrics.CancelOrderGrpcRequests.Inc()

	command := v1.NewCancelOrderCommand(req.GetAggregateID(), req.GetCancelReason())
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.Errorf("(validate) err: {%v}", err)
		tracing.TraceErr(span, err)
		return nil, s.errResponse(err)
	}

	if err := s.os.Commands.CancelOrder.Handle(ctx, command); err != nil {
		s.log.Errorf("(CancelOrder.Handle) orderID: {%s}, err: {%v}", req.GetAggregateID(), err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("(CancelOrder): AggregateID: {%s}", req.GetAggregateID())
	return &orderService.CancelOrderRes{}, nil
}

func (s *orderGrpcService) CompleteOrder(ctx context.Context, req *orderService.CompleteOrderReq) (*orderService.CompleteOrderRes, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "orderGrpcService.CompleteOrder")
	defer span.Finish()
	span.LogFields(log.String("CompleteOrder req", req.String()))
	s.metrics.CompleteOrderGrpcRequests.Inc()

	command := v1.NewCompleteOrderCommand(req.GetAggregateID(), time.Now())
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.Errorf("(validate) err: {%v}", err)
		tracing.TraceErr(span, err)
		return nil, s.errResponse(err)
	}

	if err := s.os.Commands.CompleteOrder.Handle(ctx, command); err != nil {
		s.log.Errorf("(CompleteOrder.Handle) orderID: {%s}, err: {%v}", req.GetAggregateID(), err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("(CompleteOrder): AggregateID: {%s}", req.GetAggregateID())
	return &orderService.CompleteOrderRes{}, nil
}

func (s *orderGrpcService) ChangeDeliveryAddress(ctx context.Context, req *orderService.ChangeDeliveryAddressReq) (*orderService.ChangeDeliveryAddressRes, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "orderGrpcService.ChangeDeliveryAddress")
	defer span.Finish()
	span.LogFields(log.String("ChangeDeliveryAddress req", req.String()))
	s.metrics.ChangeAddressOrderGrpcRequests.Inc()

	command := v1.NewChangeDeliveryAddressCommand(req.GetAggregateID(), req.GetDeliveryAddress())
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.Errorf("(validate) err: {%v}", err)
		tracing.TraceErr(span, err)
		return nil, s.errResponse(err)
	}

	if err := s.os.Commands.ChangeOrderDeliveryAddress.Handle(ctx, command); err != nil {
		s.log.Errorf("(ChangeOrderDeliveryAddress.Handle) orderID: {%s}, err: {%v}", req.GetAggregateID(), err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("(ChangeDeliveryAddress): AggregateID: {%s}", req.GetAggregateID())
	return &orderService.ChangeDeliveryAddressRes{}, nil
}

func (s *orderGrpcService) Search(ctx context.Context, req *orderService.SearchReq) (*orderService.SearchRes, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "orderGrpcService.Search")
	defer span.Finish()
	span.LogFields(log.String("SearchText", req.GetSearchText()), log.Int64("Page", req.GetPage()), log.Int64("Size", req.GetSize()))
	s.metrics.SearchOrderGrpcRequests.Inc()

	query := queries.NewSearchOrdersQuery(req.GetSearchText(), utils.NewPaginationQuery(int(req.GetSize()), int(req.GetPage())))
	if err := s.v.StructCtx(ctx, query); err != nil {
		s.log.Errorf("(validate) err: {%v}", err)
		tracing.TraceErr(span, err)
		return nil, s.errResponse(err)
	}

	searchResult, err := s.os.Queries.SearchOrders.Handle(ctx, query)
	if err != nil {
		s.log.Errorf("(SearchOrders.Handle) text: {%s}, err: {%v}", req.GetSearchText(), err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("(Search result): searchText: {%s}, pagination: {%+v}", req.GetSearchText(), searchResult.Pagination)
	return mappers.SearchResponseToProto(searchResult), nil
}

func (s *orderGrpcService) errResponse(err error) error {
	return grpcErrors.ErrResponse(err)
}
