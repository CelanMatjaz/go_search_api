package types

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	WithId
	DisplayName  string `json:"displayName" db:"display_name" body:""`
	Email        string `json:"email" db:"email" body:""`
	TokenVersion int    `json:"-" db:"refresh_token_version"`
	IsOauth      bool   `json:"-" db:"is_oauth"`
	WithTimestamps
	PasswordHash NullString `json:"-" db:"password_hash"`
}

type RegisterBody struct {
	DisplayName    string `json:"displayName" validate:"required,min:4,max:32"`
	Email          string `json:"email" validate:"email,required"`
	Password       string `json:"password" validate:"password,required,min:8,max:32"`
	PasswordVerify string `json:"passwordVerify"`
}

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateNewAccountData(displayName string, email string, password string) (Account, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return Account{
		DisplayName:  displayName,
		Email:        email,
		PasswordHash: NullString{sql.NullString{String: string(hash), Valid: true}},
	}, err
}

func (b LoginBody) ComparePassword(accountPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(accountPassword), []byte(b.Password))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func IsNumber(c byte) bool {
	if c >= '0' && c <= '9' {
		return true
	}

	return false
}

func IsSpecialCharacter(c byte) bool {
	if (c >= '!' && c <= '/') ||
		(c >= ':' && c <= '@') ||
		(c >= '[' && c <= '^') ||
		(c >= '{' && c <= '~') {
		return true
	}

	return false
}
