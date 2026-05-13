package employee

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	endpoints Endpoints
	logger    *logrus.Entry
}

func NewHandler(endpoints Endpoints, logger *logrus.Entry) *Handler {
	return &Handler{endpoints: endpoints, logger: logger}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	employees := router.Group("/employees")
	employees.POST("", h.Create)
	employees.GET("", h.List)
	employees.GET("/:id", h.GetByID)
	employees.PUT("/:id", h.Update)
	employees.DELETE("/:id", h.Delete)
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.endpoints.Create(c.Request.Context(), req)
	if err != nil {
		h.logger.WithError(err).Error("create employee endpoint failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	res := response.(CreateResponse)
	if res.Error != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": res.Error})
		return
	}

	c.JSON(http.StatusCreated, res.Employee)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.endpoints.GetByID(c.Request.Context(), GetByIDRequest{ID: id})
	if err != nil {
		h.logger.WithError(err).Error("get employee endpoint failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	res := response.(GetByIDResponse)
	if res.Error != "" {
		if res.Error == ErrNotFound.Error() {
			c.JSON(http.StatusNotFound, gin.H{"error": res.Error})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": res.Error})
		return
	}

	c.JSON(http.StatusOK, res.Employee)
}

func (h *Handler) List(c *gin.Context) {
	response, err := h.endpoints.List(c.Request.Context(), struct{}{})
	if err != nil {
		h.logger.WithError(err).Error("list employees endpoint failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	res := response.(ListResponse)
	if res.Error != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": res.Error})
		return
	}

	c.JSON(http.StatusOK, res.Employees)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ID = id

	response, err := h.endpoints.Update(c.Request.Context(), req)
	if err != nil {
		h.logger.WithError(err).Error("update employee endpoint failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	res := response.(UpdateResponse)
	if res.Error != "" {
		if res.Error == ErrNotFound.Error() {
			c.JSON(http.StatusNotFound, gin.H{"error": res.Error})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": res.Error})
		return
	}

	c.JSON(http.StatusOK, res.Employee)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.endpoints.Delete(c.Request.Context(), DeleteRequest{ID: id})
	if err != nil {
		h.logger.WithError(err).Error("delete employee endpoint failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	res := response.(DeleteResponse)
	if res.Error != "" {
		if res.Error == ErrNotFound.Error() {
			c.JSON(http.StatusNotFound, gin.H{"error": res.Error})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": res.Error})
		return
	}

	c.Status(http.StatusNoContent)
}

func parseID(id string) (uint, error) {
	parsed, err := strconv.ParseUint(id, 10, strconv.IntSize)
	if err != nil {
		return 0, errors.New("invalid employee id")
	}
	return uint(parsed), nil
}
