package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Handler exposes HTTP handlers for authentication.
type Handler struct {
	svc    *Service
	logger *logrus.Entry
}

// NewHandler creates a new auth Handler.
func NewHandler(svc *Service, logger *logrus.Entry) *Handler {
	return &Handler{svc: svc, logger: logger}
}

// RegisterRoutes registers auth routes on the given router group.
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/auth/login", h.Login)
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	Token string `json:"token"`
}

// Login handles POST /api/v1/auth/login.
func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}
		h.logger.WithError(err).Error("login failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, loginResponse{Token: token})
}
