package types

import "time"

type WithId struct {
	Id int `json:"id" db:"id" body:"select"`
}

type WithAccountId struct {
	AccountId int `json:"-" db:"account_id" body:"create,select"`
}

type WithTimestamps struct {
	CreatedAt time.Time `db:"created_at" json:"createdAt" body:"select"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt" body:"update,select"`
}

func (c WithId) GetId() int {
	return c.Id
}
