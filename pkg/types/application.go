package types

import (
	"strconv"
)

type Application struct {
	Common
	UserId int    `json:"user_id" db:"user_id"`
	Label  string `json:"label" db:"label"`
	Text   string `json:"text" db:"text"`
	Timestamps
}

type ApplicationSection struct {
	Common
	UserId int    `json:"user_id" db:"user_id"`
	Label  string `json:"label" db:"label"`
	Text   string `json:"text" db:"text"`
	Timestamps
}

const (
	MAX_TEXT_LENGTH = 512
)

var (
	labelError = "Value of label is not 1 to 64 characters long"
	textError  = "Value of text is not 1 to " + strconv.Itoa(MAX_TEXT_LENGTH) + " characters long"
)

func (b Application) IsValid() []string {
	errors := make([]string, 0)

	if len(b.Label) < 1 || len(b.Label) > 64 {
		errors = append(errors, labelError)
	}

	if len(b.Text) < 1 || len(b.Text) > MAX_TEXT_LENGTH {
		errors = append(errors, textError)
	}

	return errors
}

func (b ApplicationSection) IsValid() []string {
	errors := make([]string, 0)

	if len(b.Label) < 1 || len(b.Label) > 64 {
		errors = append(errors, labelError)
	}

	if len(b.Text) < 1 || len(b.Text) > MAX_TEXT_LENGTH {
		errors = append(errors, textError)
	}

	return errors
}
