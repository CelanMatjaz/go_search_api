package types

type CommonUser struct {
}

type User struct {
	Common
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Email     string `json:"email" db:"email"`
    PasswordHash string `json:"-"` 
	Timestamps
}
