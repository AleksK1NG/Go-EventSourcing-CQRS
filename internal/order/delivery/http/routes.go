package http

func (h *orderHandlers) MapRoutes() {
	h.group.POST("", h.CreateOrder())
	h.group.PUT("/pay/:id", h.PayOrder())
	h.group.PUT("/submit/:id", h.SubmitOrder())
	h.group.PUT("/:id", h.UpdateOrder())

	h.group.GET("/:id", h.GetOrderByID())
	h.group.GET("/search", h.Search())
}
