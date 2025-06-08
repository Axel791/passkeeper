package services

import (
	"fmt"
	"time"

	"github.com/Axel791/auth/internal/usecases/user/dto"

	userdomain "github.com/Axel791/auth/internal/domains/user"
	"github.com/golang-jwt/jwt/v4"
)

// TokenServiceHandler - структура сервиса работы с токеном
type TokenServiceHandler struct {
	secretKey string
}

// NewTokenService - конструктор сервиса работы с токеном
func NewTokenService(secretKey string) *TokenServiceHandler {
	return &TokenServiceHandler{secretKey: secretKey}
}

// GenerateToken - генерация токена
func (s *TokenServiceHandler) GenerateToken(claimsDTO dto.Claims) (string, error) {
	claims := jwt.MapClaims{
		"userID": claimsDTO.UserID,
		"login":  claimsDTO.Email,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.secretKey))

	if err != nil {
		return "", fmt.Errorf("error generating token: %w", err)
	}
	return signedToken, nil
}

// ValidateToken - валидация токена
func (s *TokenServiceHandler) ValidateToken(tokenStr string) (dto.Claims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return dto.Claims{}, fmt.Errorf("error parsing token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				return dto.Claims{}, fmt.Errorf("token has expired")
			}
		}

		userIDFloat, ok := claims["userID"].(float64)
		if !ok {
			return dto.Claims{}, fmt.Errorf("invalid userID in token")
		}

		email, ok := claims["email"].(string)
		if !ok {
			return dto.Claims{}, fmt.Errorf("invalid login in token")
		}

		userID, err := userdomain.NewUserID(int64(userIDFloat))
		if err != nil {
			return dto.Claims{}, fmt.Errorf("invalid user ID in token: %w", err)
		}

		return dto.Claims{
			UserID: userID,
			Email:  email,
		}, nil
	}

	return dto.Claims{}, fmt.Errorf("invalid token")
}
