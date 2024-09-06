package types

type CommonUser struct {
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Email     string `json:"email" db:"email"`
}

type User struct {
	Common
	CommonUser
	Timestamps
}

type InternalUser struct {
	User
	PasswordHash string `json:"password_hash" db:"password_hash"`
}
