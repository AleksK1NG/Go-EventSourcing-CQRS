package queries

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
}

func NewSearchOrdersQuery(searchText string) *SearchOrdersQuery {
	return &SearchOrdersQuery{SearchText: searchText}
}
