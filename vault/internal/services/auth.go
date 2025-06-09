package services

import (
	"context"
	"fmt"

	userdomain "github.com/Axel791/vault/internal/domains/user"
	"github.com/Axel791/vault/internal/services/providers"
)

type ValidateToken struct {
	validator providers.AuthValidator
}

func NewValidateToken(validator providers.AuthValidator) *ValidateToken {
	return &ValidateToken{validator: validator}
}

func (v *ValidateToken) AuthenticateToken(ctx context.Context, token string) (userdomain.UserID, error) {
	userID, err := v.validator.ValidateToken(ctx, token)
	if err != nil {
		return userdomain.UserID{}, fmt.Errorf("validate token: %w", err)
	}
	return userID, nil
}
