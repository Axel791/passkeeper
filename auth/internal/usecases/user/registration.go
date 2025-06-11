package user

import "context"

type RegistrationUseCase interface {
	Execute(ctx context.Context, email string, password string) error
}
