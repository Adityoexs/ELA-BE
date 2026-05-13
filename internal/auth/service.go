package auth

import (
	"crypto/subtle"
	"errors"
	"time"

	"github.com/Adityoexs/ELA-BE/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

// ErrInvalidCredentials is returned when the username or password does not match.
var ErrInvalidCredentials = errors.New("invalid credentials")

// Claims is the set of JWT claims embedded in every issued token.
// It deliberately includes a Role field so RBAC can be layered on top.
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// mvpUser holds an in-memory user entry used for the MVP.
// WARNING: This is a temporary placeholder for development/demo purposes only.
// Passwords are stored as plain strings solely because this map is never persisted
// and the credentials are not real secrets. Replace this with a real database
// lookup (with properly hashed passwords) before any production deployment.
type mvpUser struct {
	password string
	role     string
}

// mvpUsers is the in-memory user store used for the MVP.
// Replace (or extend) this with a real database lookup once a users table exists.
var mvpUsers = map[string]mvpUser{
	"admin@example.com": {password: "admin123", role: "admin"},
}

// Service handles JWT issuance and validation.
type Service struct {
	cfg config.JWTConfig
}

// NewService creates a new auth Service.
func NewService(cfg config.JWTConfig) *Service {
	return &Service{cfg: cfg}
}

// dummyPassword is used to perform a constant-time comparison when a user is
// not found, so that the response time does not reveal whether an email exists.
var dummyPassword = []byte("dummy-password-for-timing-safety")

// Login verifies credentials and returns a signed JWT on success.
// Password comparison uses constant-time equality to prevent timing attacks.
func (s *Service) Login(email, password string) (string, error) {
	user, ok := mvpUsers[email]
	if !ok {
		// Always perform a comparison even when the user is not found to prevent
		// user-enumeration via response-time differences.
		subtle.ConstantTimeCompare([]byte(password), dummyPassword)
		return "", ErrInvalidCredentials
	}

	if subtle.ConstantTimeCompare([]byte(password), []byte(user.password)) != 1 {
		return "", ErrInvalidCredentials
	}

	expiry := time.Duration(s.cfg.ExpirySeconds) * time.Second
	claims := Claims{
		UserID: email,
		Email:  email,
		Role:   user.role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.Secret))
}

// ValidateToken parses and validates a JWT string, returning the embedded claims.
func (s *Service) ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.cfg.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
