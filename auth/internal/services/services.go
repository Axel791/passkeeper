package services

import "github.com/Axel791/auth/interanal/usecases/user/dto"

type HashPasswordService interface {
	Hash(string) string
}

type TokenService interface {
	GenerateToken(claimsDTO dto.Claims) (string, error)
	ValidateToken(tokenStr string) (dto.Claims, error)
}
