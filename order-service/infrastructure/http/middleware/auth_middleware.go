package middleware

import (
	"net/http"
	"strings"

	authclient "github.com/adopabianko/commerce/order-service/infrastructure/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(ac *authclient.GRPCAuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			return
		}

		token := parts[1]
		userID, err := ac.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// set user_id ke context
		c.Set("user_id", userID)
		c.Next()
	}
}
