package types

type ApplicationSection struct {
	Common
	AccountId int    `json:"-" db:"account_id" body:""`
	Label     string `json:"label" db:"label" validate:"required,min:1,max:32" body:""`
	Text      string `json:"text" db:"text" validate:"required,min:1,max:1024" body:""`
	WithTags
	Timestamps
}

type ApplicationPreset struct {
	Common
	AccountId int    `json:"-" db:"account_id" body:""`
	Label     string `json:"label" db:"label" validate:"required,min:1,max:32" body:""`
	WithTags
	Timestamps
}
