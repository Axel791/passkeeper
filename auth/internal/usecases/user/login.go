package user

import "context"

type Login interface {
	Execute(
		ctx context.Context,
		email string,
		password string,
	) (string, error)
}
