package v1

func (h *orderHandlers) MapRoutes() {
	h.group.POST("", h.CreateOrder())
	h.group.PUT("/pay/:id", h.PayOrder())
	h.group.PUT("/submit/:id", h.SubmitOrder())
	h.group.PUT("/cart/:id", h.UpdateShoppingCart())
	h.group.POST("/cancel/:id", h.CancelOrder())
	h.group.POST("/complete/:id", h.CompleteOrder())
	h.group.PUT("/address/:id", h.ChangeDeliveryAddress())

	h.group.GET("/:id", h.GetOrderByID())
	h.group.GET("/search", h.Search())
}
