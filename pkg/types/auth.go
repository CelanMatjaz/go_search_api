package types

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	Common
	DisplayName  string         `json:"displayName" db:"display_name"`
	Email        string         `json:"email" db:"email"`
	PasswordHash sql.NullString `json:"-" db:"password_hash"`
	TokenVersion int            `json:"-" db:"refresh_token_version"`
	Timestamps
}

type RegisterBody struct {
	DisplayName    string `json:"displayName"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	PasswordVerify string `json:"passwordVerify"`
}

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateNewAccountData(displayName string, email string, password string) (*Account, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return &Account{
		DisplayName:  displayName,
		Email:        email,
		PasswordHash: sql.NullString{String: string(hash), Valid: true},
	}, err
}

func (b RegisterBody) Verify() []string {
	errors := make([]string, 0)

	if b.DisplayName == "" {
		errors = append(errors, "Property displayName missing from JSON body")
	}
	if b.Email == "" {
		errors = append(errors, "Property email missing from JSON body")
	}
	if b.Password == "" {
		errors = append(errors, "Property password missing from JSON body")
	}
	if b.PasswordVerify == "" {
		errors = append(errors, "Property passwordVerify missing from JSON body")
	}

	if len(errors) > 0 {
		return errors
	}

	number := false
	specialCharacter := false
	upperCase := false
	lowerCase := false

	for _, c := range []byte(b.Password) {
		if IsNumber(c) {
			number = true
		} else if IsSpecialCharacter(c) {
			specialCharacter = true
		} else if c >= 'a' && c <= 'z' {
			lowerCase = true
		} else if c >= 'A' && c <= 'Z' {
			upperCase = true
		}
	}

	if !number {
		errors = append(errors, "Password requires at least one number")
	}
	if !specialCharacter {
		errors = append(errors, "Password requires at least one special character")
	}
	if !upperCase {
		errors = append(errors, "Password requires at least upper case letter")
	}
	if !lowerCase {
		errors = append(errors, "Password requires at least lower case letter")
	}
	if b.Password != b.PasswordVerify {
		errors = append(errors, "Passwords do not match")
	}
	if len(b.Password) < 8 && len(b.Password) > 32 {
		errors = append(errors, "Password length is not between 8 and 32 characters long")
	}

	return errors
}

func (b LoginBody) Verify() []string {
	errors := make([]string, 0)

	if b.Email == "" {
		errors = append(errors, "Property email missing from JSON body")
	}
	if b.Password == "" {
		errors = append(errors, "Property password missing from JSON body")
	}

	return errors
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
