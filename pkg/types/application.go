package types

type Application struct {
	Common
	UserId int    `json:"user_id" db:"user_id"`
	Label  string `json:"name" db:"label"`
	Text   string `json:"text" db:"text"`
	Timestamps
}

type ApplicationSection struct {
	Common
	UserId int    `json:"user_id" db:"user_id"`
	Label  string `json:"name" db:"label"`
	Text   string `json:"text" db:"text"`
	Timestamps
}

var (
	labelError = "Value of label is not 1 to 64 characters long"
	textError  = "Value of text is not 1 to 64 characters long"
)

func (b Application) IsValid() []string {
	errors := make([]string, 0)

	if len(b.Label) < 1 || len(b.Label) > 64 {
		errors = append(errors, labelError)
	}

	if len(b.Text) < 1 || len(b.Text) > 64 {
		errors = append(errors, textError)
	}

	return errors
}
