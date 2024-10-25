package types

type ScanTag struct {
	Id        NullInt64  `db:"id"`
	AccountId NullInt64  `db:"account_id"`
	Label     NullString `db:"label"`
	Color     NullString `db:"color"`
}

func (t ScanTag) Tag() Tag {
	return Tag{
		WithId: WithId{
			Id: int(t.Id.Int64),
		},
		WithAccountId: WithAccountId{
			AccountId: int(t.AccountId.Int64),
		},
		Label: t.Label.String,
		Color: t.Color.String,
	}
}

type Tag struct {
	WithId
	WithAccountId
	Label string `json:"label" db:"label" body:"" validate:"required,min:1,max:32"`
	Color string `json:"color" db:"color" body:"" validate:"required,len:7"`
}

func CreateTag(accountId int, label string, color string) Tag {
	return Tag{
		WithAccountId: WithAccountId{
			AccountId: accountId,
		},
		Label: label,
		Color: color,
	}
}

type TagBody struct {
	Label string `json:"label" db:"label" body:"" validate:"required,min:1,max:32"`
	Color string `json:"color" db:"color" body:"" validate:"required,len:7"`
}

type WithTags struct {
	TagIds []int `json:"tagIds,omitempty" body:""`
}

type RecordWithTags[T any] struct {
	Record T     `json:"record"`
	Tags   []Tag `json:"tags,omitempty"`
}

func (wt WithTags) GetTagCount() int {
	return len(wt.TagIds)
}
