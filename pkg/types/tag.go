package types

type Tag struct {
	Common
	AccountId int    `json:"accountId" db:"account_id"`
	Label     string `json:"label" db:"label"`
	Color     string `json:"color" db:"color"`
}

type TagBody struct {
	Label string `json:"label" db:"label"`
	Color string `json:"color" db:"color"`
}

func (b TagBody) Verify() []string {
	errors := make([]string, 0)

	if b.Label == "" {
		errors = append(errors, "Property label missing from JSON body")
	}
	if b.Color == "" {
		errors = append(errors, "Property color missing from JSON body")
	}

	return errors
}
