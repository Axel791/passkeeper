package config

import (
	"github.com/Axel791/auth/internal/services"
	groupusecases "github.com/Axel791/auth/internal/usecases/group"
	groupscenario "github.com/Axel791/auth/internal/usecases/group/scenarios"
	userusecases "github.com/Axel791/auth/internal/usecases/user"
	userscenario "github.com/Axel791/auth/internal/usecases/user/scenarios"
	log "github.com/sirupsen/logrus"
)

func newValidate(
	providers Providers,
	logger *log.Logger,
	tokenService services.TokenService,
) userusecases.Validate {
	validate := userscenario.NewValidate(
		logger,
		tokenService,
		providers.UserProvider,
	)
	return validate
}

func newLogin(
	logger *log.Logger,
	providers Providers,
	hashService services.HashPasswordService,
	tokenService services.TokenService,
) userusecases.Login {
	login := userscenario.NewLoginScenario(
		logger,
		providers.UserProvider,
		hashService,
		tokenService,
	)
	return login
}

func newRegister(
	logger *log.Logger,
	providers Providers,
	hashService services.HashPasswordService,
) userusecases.RegistrationUseCase {
	register := userscenario.NewRegistrationScenario(logger, providers.UserProvider, hashService)
	return register
}

func newGetUserGroup(logger *log.Logger, providers Providers) groupusecases.GetGroupByUserID {
	group := groupscenario.NewGetGroupByUserID(logger, providers.GroupProvider)
	return group
}

type UseCases struct {
	Validate      userusecases.Validate
	Login         userusecases.Login
	Registration  userusecases.RegistrationUseCase
	GroupUseCases groupusecases.GetGroupByUserID
}

func NewUseCases(
	logger *log.Logger,
	providers Providers,
	hashService services.HashPasswordService,
	tokenService services.TokenService,
) *UseCases {
	return &UseCases{
		Validate:      newValidate(providers, logger, tokenService),
		Login:         newLogin(logger, providers, hashService, tokenService),
		Registration:  newRegister(logger, providers, hashService),
		GroupUseCases: newGetUserGroup(logger, providers),
	}
}
