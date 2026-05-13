package auth_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Adityoexs/ELA-BE/internal/auth"
	"github.com/Adityoexs/ELA-BE/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func testConfig() config.JWTConfig {
	return config.JWTConfig{
		Secret:        "test-secret",
		ExpirySeconds: 3600,
	}
}

func newTestRouter(svc *auth.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	log := logrus.NewEntry(logrus.New())
	h := auth.NewHandler(svc, log)
	v1 := r.Group("/api/v1")
	h.RegisterRoutes(v1)
	protected := v1.Group("", auth.Middleware(svc))
	protected.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"pong": true})
	})
	return r
}

// --- Login tests ---

func TestLogin_Success(t *testing.T) {
	svc := auth.NewService(testConfig())
	r := newTestRouter(svc)

	body := `{"email":"admin@example.com","password":"admin123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["token"] == "" {
		t.Error("expected non-empty token in response")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	svc := auth.NewService(testConfig())
	r := newTestRouter(svc)

	body := `{"email":"admin@example.com","password":"wrong"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestLogin_UnknownUser(t *testing.T) {
	svc := auth.NewService(testConfig())
	r := newTestRouter(svc)

	body := `{"email":"nobody@example.com","password":"whatever"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestLogin_MissingFields(t *testing.T) {
	svc := auth.NewService(testConfig())
	r := newTestRouter(svc)

	body := `{"email":"admin@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

// --- Middleware tests ---

func TestMiddleware_NoHeader(t *testing.T) {
	svc := auth.NewService(testConfig())
	r := newTestRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestMiddleware_InvalidToken(t *testing.T) {
	svc := auth.NewService(testConfig())
	r := newTestRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)
	req.Header.Set("Authorization", "Bearer not.a.valid.token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestMiddleware_WrongScheme(t *testing.T) {
	svc := auth.NewService(testConfig())
	r := newTestRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)
	req.Header.Set("Authorization", "Basic dXNlcjpwYXNz")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestMiddleware_ValidToken(t *testing.T) {
	svc := auth.NewService(testConfig())
	r := newTestRouter(svc)

	// First obtain a real token via login
	body := `{"email":"admin@example.com","password":"admin123"}`
	loginReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(body))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	r.ServeHTTP(loginW, loginReq)

	var loginResp map[string]string
	if err := json.Unmarshal(loginW.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("failed to decode login response: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)
	req.Header.Set("Authorization", "Bearer "+loginResp["token"])
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestMiddleware_ExpiredToken(t *testing.T) {
	// Create a service with 0-second expiry to force immediate expiration
	expiredCfg := config.JWTConfig{Secret: "test-secret", ExpirySeconds: -1}
	svc := auth.NewService(expiredCfg)

	token, err := svc.Login("admin@example.com", "admin123")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	// Validate that an expired token is rejected
	_, err = svc.ValidateToken(token)
	if err == nil {
		t.Error("expected error for expired token, got nil")
	}
}
