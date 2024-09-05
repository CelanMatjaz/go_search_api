package types

import "time"

type Common struct {
	Id int `json:"id" db:"id"`
}

type Timestamps struct {
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

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

type Resume struct {
	Common
	Name string `json:"name" db:"name"`
	Note string `json:"note" db:"note"`
	Timestamps
}

type JobListing struct {
	Common
	Url     string `json:"url" db:"url"`
	Company string `json:"company" db:"company"`
}

type ResumeTag struct {
	Common
	Label string `json:"label" db:"label"`
}
