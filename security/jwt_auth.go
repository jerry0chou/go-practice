package security

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the claims structure for JWT tokens
type JWTClaims struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// JWTAuth handles JWT token operations
type JWTAuth struct {
	secretKey []byte
}

// NewJWTAuth creates a new JWT authentication instance
func NewJWTAuth(secretKey string) *JWTAuth {
	return &JWTAuth{
		secretKey: []byte(secretKey),
	}
}

// GenerateToken creates a new JWT token for a user
func (j *JWTAuth) GenerateToken(userID, username string, roles []string, expirationHours int) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expirationHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-practice-app",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// ValidateToken validates and parses a JWT token
func (j *JWTAuth) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken generates a new token with extended expiration
func (j *JWTAuth) RefreshToken(tokenString string, expirationHours int) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Generate new token with same user info but new expiration
	return j.GenerateToken(claims.UserID, claims.Username, claims.Roles, expirationHours)
}

// ExtractUserInfo extracts user information from a valid token
func (j *JWTAuth) ExtractUserInfo(tokenString string) (userID, username string, roles []string, err error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", "", nil, err
	}

	return claims.UserID, claims.Username, claims.Roles, nil
}
