package types

type ResumeSection struct {
	WithId
	WithAccountId
	Label string `json:"label" db:"label" validate:"required,min:1,max:32" body:"create,update"`
	Text  string `json:"text" db:"text" validate:"required,min:1,max:1024" body:"create,update"`
	*WithTags
	WithTimestamps
}

type ResumePreset struct {
	WithId
	WithAccountId
	Label string `json:"label" db:"label" validate:"required,min:1,max:32" body:"create,update"`
	*WithTags
	WithTimestamps
}
