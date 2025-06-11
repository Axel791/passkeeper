package scenarios

import (
	"context"
	"fmt"

	"github.com/Axel791/auth/internal/services"
	"github.com/Axel791/auth/internal/usecases/user/dto"
	"github.com/Axel791/auth/internal/usecases/user/providers"
	log "github.com/sirupsen/logrus"
)

type Validate struct {
	logger         *log.Logger
	tokenService   services.TokenService
	userRepository providers.UserRepository
}

func NewValidate(
	logger *log.Logger,
	tokenService services.TokenService,
	userRepository providers.UserRepository,
) *Validate {
	return &Validate{
		logger:         logger,
		tokenService:   tokenService,
		userRepository: userRepository,
	}
}

func (v *Validate) Execute(ctx context.Context, token string) (dto.User, error) {
	claims, err := v.tokenService.ValidateToken(token)
	if err != nil {
		return dto.User{}, fmt.Errorf("invalid token: %w", err)
	}

	user, err := v.userRepository.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return dto.User{}, fmt.Errorf("invalid user: %w", err)
	}

	return dto.User{
		ID:    user.ID().ToInt64(),
		Email: user.Email(),
	}, nil
}
