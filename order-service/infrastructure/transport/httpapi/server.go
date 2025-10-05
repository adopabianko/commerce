package httpapi

import (
	"net/http"

	"github.com/adopabianko/commerce/order-service/infrastructure/auth"
	"github.com/adopabianko/commerce/order-service/infrastructure/http/middleware"
	"github.com/adopabianko/commerce/order-service/internal/usecase"
	"github.com/gin-gonic/gin"
)

type Handler struct{ Place *usecase.PlaceOrder }

func New(place *usecase.PlaceOrder) *Handler { return &Handler{Place: place} }

func (h *Handler) Routes(r *gin.Engine, authClient *auth.GRPCAuthClient) {
	r.Use(middleware.AuthMiddleware(authClient))
	r.POST("/orders", h.placeOrder)
}

func (h *Handler) placeOrder(c *gin.Context) {
	var req struct {
		Items []struct {
			SKU string `json:"sku"`
			Qty int32  `json:"qty"`
		} `json:"items"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	items := make([]usecase.Item, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, usecase.Item{SKU: it.SKU, Qty: it.Qty})
	}
	resp, err := h.Place.Exec(c, usecase.Request{Items: items})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
