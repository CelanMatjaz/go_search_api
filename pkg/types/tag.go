package types

import (
	"database/sql"
	"encoding/json"
)

type NullString struct {
	sql.NullString
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(ns.String)
}

type NullInt64 struct{ sql.NullInt64 }

func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(ni.Int64)
}

type Tag struct {
	Id        NullInt64  `json:"id" db:"id"`
    AccountId NullInt64  `json:"-" db:"account_id"`
	Label     NullString `json:"label" db:"label"`
	Color     NullString `json:"color" db:"color"`
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
