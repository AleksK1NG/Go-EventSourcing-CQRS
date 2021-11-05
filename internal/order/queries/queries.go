package queries

type OrderQueries struct {
	GetOrderByIDQuery GetOrderByIDQueryHandler
}

func NewOrderQueries(getOrderByIDQuery GetOrderByIDQueryHandler) *OrderQueries {
	return &OrderQueries{GetOrderByIDQuery: getOrderByIDQuery}
}

type GetOrderByIDQuery struct {
	ID string
}

func NewGetOrderByIDQuery(ID string) *GetOrderByIDQuery {
	return &GetOrderByIDQuery{ID: ID}
}
