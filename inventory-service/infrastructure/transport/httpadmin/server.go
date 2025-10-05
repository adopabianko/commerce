package httpadmin

import (
	"net/http"

	domain "github.com/adopabianko/commerce/inventory-service/internal/domain/inventory"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine, repo domain.Repository) {
	r.POST("/seed", func(c *gin.Context) {
		var req struct {
			Products []domain.Product `json:"products"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := repo.BulkUpsertProducts(c, req.Products); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true, "count": len(req.Products)})
	})
}
