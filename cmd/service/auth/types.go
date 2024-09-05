package auth

import (
	"errors"

	"github.com/CelanMatjaz/job_application_tracker_api/cmd/types"
)

var InvalidBodyErr = errors.New("Provided JSON body is not valid")
var PasswordsDoNotMatchErr = errors.New("Passwords do not match")

type RegisterBody struct {
	FirstName      *string `json:"first_name"`
	LastName       *string `json:"last_name"`
	Email          *string `json:"email"`
	Password       *string `json:"password"`
	PasswordVerify *string `json:"password_verify"`
}

func (r *RegisterBody) IsValid() error {
	if r.FirstName == nil ||
		r.LastName == nil ||
		r.Email == nil ||
		r.Password == nil ||
		r.PasswordVerify == nil {
		return InvalidBodyErr
	}

	if *r.Password != *r.PasswordVerify {
		return PasswordsDoNotMatchErr
	}

	return nil
}

func (r *RegisterBody) CreateInternalUser(passwordHash string) types.InternalUser {
	return types.InternalUser{
		User: types.User{
			CommonUser: types.CommonUser{
				FirstName: *r.FirstName,
				LastName:  *r.LastName,
				Email:     *r.Email,
			},
		},
		PasswordHash: passwordHash,
	}
}

type LoginBody struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func (l *LoginBody) IsValid() error {
	if l.Email == nil || l.Password == nil {
		return InvalidBodyErr
	}

	return nil
}
