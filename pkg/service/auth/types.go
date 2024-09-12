package auth

import (
	"database/sql"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

type RegisterBody struct {
	DisplayName    string `json:"display_name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	PasswordVerify string `json:"password_verify"`
}

func (r *RegisterBody) IsValid() []string {
	errors := make([]string, 0)

	if r.DisplayName == "" ||
		r.Email == "" ||
		r.Password == "" ||
		r.PasswordVerify == "" {
		errors = append(errors, "At least 1 required field is missing")
	}

	specialCharacterOrNumber := false
	upperCase := false
	lowerCase := false
	for _, c := range []byte(r.Password) {
		if isSpecialCharacterOrNumber(c) {
			specialCharacterOrNumber = true
		} else if c >= 'a' && c <= 'z' {
			lowerCase = true
		} else if c >= 'A' && c <= 'Z' {
			upperCase = true
		}
	}

	if !specialCharacterOrNumber {
		errors = append(errors, "Password requires at least one special character or number")
	}
	if !upperCase {
		errors = append(errors, "Password requires at least upper case letter")
	}
	if !lowerCase {
		errors = append(errors, "Password requires at least lower case letter")
	}

	if r.Password != r.PasswordVerify {
		errors = append(errors, "Passwords do not match")
	}

	return errors
}

func (r *RegisterBody) CreateUser(passwordHash string) types.User {
	return types.User{
		DisplayName:  r.DisplayName,
		Email:        r.Email,
		PasswordHash: sql.NullString{String: passwordHash, Valid: true},
		TokenVersion: 1,
	}
}

type LoginBody struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func (l *LoginBody) IsValid() error {
	if l.Email == nil || l.Password == nil {
		return types.InvalidBodyErr
	}

	return nil
}

func isSpecialCharacterOrNumber(c byte) bool {
	if (c >= '!' && c <= '@') || (c >= '[' && c <= '^') {
		return true
	}

	return false
}
