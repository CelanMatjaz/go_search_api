package types

type Resume struct {
	Common
	UserId int    `json:"user_id" db:"user_id"`
	Name   string `json:"name" db:"name"`
	Note   string `json:"note" db:"note"`
	Timestamps
}

type ResumeSection struct {
	Common
	UserId   int    `json:"user_id" db:"user_id"`
	Label    string `json:"name" db:"label"`
	Markdown string `json:"markdown" db:"markdown"`
	Timestamps
}

