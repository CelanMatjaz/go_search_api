package types

import "time"

type WithId struct {
	Id int `json:"id" db:"id"`
}

type WithAccountId struct {
	AccountId int `json:"-" db:"account_id" body:"omit"`
}

type Timestamps struct {
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

func (c WithId) GetId() int {
	return c.Id
}
