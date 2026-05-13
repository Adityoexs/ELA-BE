package http

import (
	"net/http"

	"github.com/Adityoexs/ELA-BE/internal/auth"
	"github.com/Adityoexs/ELA-BE/internal/employee"
	"github.com/gin-gonic/gin"
)

func NewRouter(employeeHandler *employee.Handler, authHandler *auth.Handler, authSvc *auth.Service) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := router.Group("/api/v1")

	// Public auth routes (no JWT required)
	authHandler.RegisterRoutes(v1)

	// Protected employee routes (JWT required)
	protected := v1.Group("", auth.Middleware(authSvc))
	employeeHandler.RegisterRoutes(protected)

	return router
}
