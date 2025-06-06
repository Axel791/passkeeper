package dto

import userdomain "github.com/Axel791/auth/interanal/domains/user"

type Claims struct {
	UserID userdomain.UserID
	Email  string
}
