package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// HashPasswordServiceHandler - структура сервиса хэширования паролей
type HashPasswordServiceHandler struct {
	secretKey string
}

// NewHashPasswordService - конструктор сервиса хэширования паролей
func NewHashPasswordService(secretKey string) *HashPasswordServiceHandler {
	return &HashPasswordServiceHandler{secretKey}
}

// Hash - хэшируем пароль
func (s *HashPasswordServiceHandler) Hash(password string) string {
	h := hmac.New(sha256.New, []byte(s.secretKey))
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}
