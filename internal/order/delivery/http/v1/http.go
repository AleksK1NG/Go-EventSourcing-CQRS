package v1

import "github.com/labstack/echo/v4"

type OrderHandlers interface {
	CreateOrder() echo.HandlerFunc
	PayOrder() echo.HandlerFunc
	SubmitOrder() echo.HandlerFunc
	UpdateShoppingCart() echo.HandlerFunc

	GetOrderByID() echo.HandlerFunc
	Search() echo.HandlerFunc
}
