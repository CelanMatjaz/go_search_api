package types

type ResumeSection struct {
    WithId
	WithAccountId
	Label string `json:"label" db:"label" validate:"required,min:1,max:32" body:""`
	Text  string `json:"text" db:"text" validate:"required,min:1,max:1024" body:""`
	*WithTags
	Timestamps
}

type ResumePreset struct {
    WithId
	WithAccountId
	Label     string `json:"label" db:"label" validate:"required,min:1,max:32" body:""`
	*WithTags
	Timestamps
}
