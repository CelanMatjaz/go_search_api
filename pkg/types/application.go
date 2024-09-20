package types

type ApplicationSection struct {
	Common
	AccountId int    `json:"accountId" db:"account_id"`
	Label     string `json:"label" db:"label"`
	Text      string `json:"text" db:"text"`
	Timestamps
}

type ApplicationPreset struct {
	Common
	AccountId int    `json:"accountId" db:"account_id"`
	Label     string `json:"label" db:"label"`
	Timestamps
}

type ApplicationPresetBody struct {
	Label      string `json:"label" db:"label"`
	SectionIds []int  `json:"sectionIds"`
}

func (b ApplicationPresetBody) Verify() []string {
	errors := make([]string, 0)

	if b.Label == "" {
		errors = append(errors, "Property label missing from JSON body")
	}

	if len(b.Label) > 64 {
		errors = append(errors, "Label cannot be more than 64 characters long")
	}

	return errors
}

type ApplicationSectionBody struct {
	Label string `json:"label" db:"label"`
	Text  string `json:"text" db:"text"`
}

func (b ApplicationSectionBody) Verify() []string {
	errors := make([]string, 0)

	if b.Label == "" {
		errors = append(errors, "Property label missing from JSON body")
	}
	if b.Text == "" {
		errors = append(errors, "Property text missing from JSON body")
	}

	if len(errors) > 0 {
		return errors
	}

	if len(b.Label) > 64 {
		errors = append(errors, "Label cannot be more than 64 characters long")
	}
	if len(b.Text) > 64 {
		errors = append(errors, "Text cannot be more than 1024 characters long")
	}

	return errors
}
