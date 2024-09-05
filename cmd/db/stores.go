package db

import "github.com/CelanMatjaz/job_application_tracker_api/cmd/types"

type AuthStore interface {
	GetInternalUserById(id int) (types.InternalUser, error)
	GetInternalUserByEmail(email string) (types.InternalUser, error)
	CreateUser(user types.InternalUser) (types.InternalUser, error)
}
