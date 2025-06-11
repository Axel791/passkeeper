package scenarios

import (
	"context"

	"github.com/Axel791/appkit"
	userdomain "github.com/Axel791/auth/internal/domains/user"
	"github.com/Axel791/auth/internal/services"
	"github.com/Axel791/auth/internal/usecases/user/providers"
	log "github.com/sirupsen/logrus"
)

type RegistrationScenario struct {
	logger              *log.Logger
	userRepository      providers.UserRepository
	hashPasswordService services.HashPasswordService
}

func NewRegistrationScenario(
	logger *log.Logger,
	userRepository providers.UserRepository,
	hashPasswordService services.HashPasswordService,
) *RegistrationScenario {
	return &RegistrationScenario{
		logger:              logger,
		userRepository:      userRepository,
		hashPasswordService: hashPasswordService,
	}
}

func (s *RegistrationScenario) Execute(
	ctx context.Context,
	email string,
	password string,
) error {
	hashedPassword := s.hashPasswordService.Hash(password)
	user := userdomain.NewCreateUser(
		email,
		hashedPassword,
	)

	err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		s.logger.WithError(err).Error("err to create user")

		return appkit.Wrap(appkit.Unknown, "err creating user", err)
	}

	return nil
}
