// package http

// import (
// 	"fmt"
// 	"net/http"
// 	"os"

// 	"github.com/adopabianko/commerce/user-service/internal/repository"
// 	"github.com/adopabianko/commerce/user-service/internal/service"
// 	"github.com/gin-gonic/gin"
// 	"gorm.io/gorm"
// )

// func RunHTTPServer(db *gorm.DB) error {
// 	repo := repository.NewUserRepository(db)
// 	jwtSecret := os.Getenv("JWT_SECRET")
// 	if jwtSecret == "" {
// 		jwtSecret = "secret"
// 	}
// 	svc := service.NewAuthService(repo, jwtSecret)

// 	r := gin.Default()

// 	r.POST("/register", func(c *gin.Context) {
// 		var req struct{ Email, Password, Name string }
// 		if err := c.BindJSON(&req); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		u, err := svc.Register(req.Email, req.Password, req.Name)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusOK, gin.H{"id": u.ID, "email": u.Email, "name": u.Name})
// 	})

// 	r.POST("/login", func(c *gin.Context) {
// 		var req struct{ Email, Password string }
// 		if err := c.BindJSON(&req); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		token, err := svc.Login(req.Email, req.Password)
// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusOK, gin.H{"token": token})
// 	})

// 	port := os.Getenv("HTTP_PORT")
// 	if port == "" {
// 		port = "8080"
// 	}
// 	addr := fmt.Sprintf(":%s", port)
// 	return r.Run(addr)
// }

package httpsvr

import (
	"net/http"

	domain "github.com/adopabianko/commerce/user-service/internal/domain/user"
	"github.com/adopabianko/commerce/user-service/internal/usecase"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine, svc *usecase.Service, repo domain.Repository, jwtSecret string) {
	r.POST("/register", func(c *gin.Context) {
		var req struct{ Email, Password, Name string }
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		u, err := svc.Register(req.Email, req.Password, req.Name)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": u.ID, "email": u.Email, "name": u.Name})
	})

	r.POST("/login", func(c *gin.Context) {
		var req struct{ Email, Password string }
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		token, err := svc.Login(req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	})
}
