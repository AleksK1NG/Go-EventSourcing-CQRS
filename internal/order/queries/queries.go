package queries

import "github.com/AleksK1NG/es-microservice/pkg/utils"

type OrderQueries struct {
	GetOrderByID GetOrderByIDQueryHandler
	SearchOrders SearchOrdersQueryHandler
}

func NewOrderQueries(getOrderByID GetOrderByIDQueryHandler, searchOrders SearchOrdersQueryHandler) *OrderQueries {
	return &OrderQueries{GetOrderByID: getOrderByID, SearchOrders: searchOrders}
}

type GetOrderByIDQuery struct {
	ID string
}

func NewGetOrderByIDQuery(ID string) *GetOrderByIDQuery {
	return &GetOrderByIDQuery{ID: ID}
}

type SearchOrdersQuery struct {
	SearchText string `json:"searchText"`
	Pq         *utils.Pagination
}

func NewSearchOrdersQuery(searchText string, pq *utils.Pagination) *SearchOrdersQuery {
	return &SearchOrdersQuery{SearchText: searchText, Pq: pq}
}
