package service

import (
	"github.com/AleksK1NG/es-microservice/config"
	"github.com/AleksK1NG/es-microservice/internal/order/commands/v1"
	"github.com/AleksK1NG/es-microservice/internal/order/queries"
	"github.com/AleksK1NG/es-microservice/internal/order/repository"
	"github.com/AleksK1NG/es-microservice/pkg/es"
	"github.com/AleksK1NG/es-microservice/pkg/logger"
)

type OrderService struct {
	Commands *v1.OrderCommands
	Queries  *queries.OrderQueries
}

func NewOrderService(
	log logger.Logger,
	cfg *config.Config,
	es es.AggregateStore,
	mongoRepo repository.OrderMongoRepository,
	elasticRepository repository.ElasticOrderRepository,
) *OrderService {

	createOrderHandler := v1.NewCreateOrderHandler(log, cfg, es)
	orderPaidHandler := v1.NewOrderPaidHandler(log, cfg, es)
	submitOrderHandler := v1.NewSubmitOrderHandler(log, cfg, es)
	updateOrderCmdHandler := v1.NewUpdateShoppingCartCmdHandler(log, cfg, es)
	cancelOrderCommandHandler := v1.NewCancelOrderCommandHandler(log, cfg, es)
	deliveryOrderCommandHandler := v1.NewCompleteOrderCommandHandler(log, cfg, es)
	changeOrderDeliveryAddressCmdHandler := v1.NewChangeDeliveryAddressCmdHandler(log, cfg, es)

	getOrderByIDHandler := queries.NewGetOrderByIDHandler(log, cfg, es, mongoRepo)
	searchOrdersHandler := queries.NewSearchOrdersHandler(log, cfg, es, elasticRepository)

	orderCommands := v1.NewOrderCommands(
		createOrderHandler,
		orderPaidHandler,
		submitOrderHandler,
		updateOrderCmdHandler,
		cancelOrderCommandHandler,
		deliveryOrderCommandHandler,
		changeOrderDeliveryAddressCmdHandler,
	)
	orderQueries := queries.NewOrderQueries(getOrderByIDHandler, searchOrdersHandler)

	return &OrderService{Commands: orderCommands, Queries: orderQueries}
}
