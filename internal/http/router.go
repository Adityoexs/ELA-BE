package http

import (
	"net/http"

	"github.com/Adityoexs/ELA-BE/internal/employee"
	"github.com/gin-gonic/gin"
)

func NewRouter(employeeHandler *employee.Handler) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := router.Group("/api/v1")
	employeeHandler.RegisterRoutes(v1)

	return router
}
