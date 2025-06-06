package scenarios

import (
	"context"
	"errors"
	"github.com/Axel791/appkit"
	"github.com/Axel791/auth/internal/services"
	"github.com/Axel791/auth/internal/usecases/user/dto"
	"github.com/Axel791/auth/internal/usecases/user/providers"
	log "github.com/sirupsen/logrus"
)

type LoginScenario struct {
	logger              *log.Logger
	userRepository      providers.UserRepository
	hashPasswordService services.HashPasswordService
	tokenService        services.TokenService
}

func NewLoginScenario(
	logger *log.Logger,
	userRepository providers.UserRepository,
	hashPasswordService services.HashPasswordService,
	tokenService services.TokenService,
) *LoginScenario {
	return &LoginScenario{
		logger:              logger,
		userRepository:      userRepository,
		hashPasswordService: hashPasswordService,
		tokenService:        tokenService,
	}
}

func (s *LoginScenario) Execute(
	ctx context.Context,
	email string,
	password string,
) (string, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, providers.ErrUserNotFound) {
			return "", appkit.New(appkit.Unauthorized, "user not found")
		}

		s.logger.WithError(err).Error("error retrieving user")

		return "", appkit.Wrap(appkit.Unknown, "unknown error", err)
	}

	hashedPassword := s.hashPasswordService.Hash(password)

	if hashedPassword != user.PwHash() {
		return "", appkit.New(appkit.Unauthorized, "invalid password")
	}

	claims := dto.Claims{
		UserID: user.ID(),
		Email:  user.Email(),
	}

	token, err := s.tokenService.GenerateToken(claims)
	if err != nil {
		s.logger.WithError(err).Error("error generating token")

		return "", appkit.Wrap(appkit.Unknown, "error generating token", err)
	}
	return token, nil
}
