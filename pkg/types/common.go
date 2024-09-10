package types

type Tag struct {
	Common
	UserId int    `json:"user_id" db:"user_id"`
	Label  string `json:"label" db:"label"`
}
