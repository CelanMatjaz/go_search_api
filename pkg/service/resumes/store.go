package resumes

import (
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

func NewStore(connection *db.DbConnection) *db.GenericStore[types.Resume] {
	tableName := "resumes"
	scanFields := []string{"id", "user_id", "name", "note", "created_at", "updated_at"}
	neededFields := []string{"user_id", "name", "note"}
	updateFields := []string{"name", "note"}
	whereClause := "WHERE id = $1 AND user_id = $2"

	return &db.GenericStore[types.Resume]{
		Db:              connection.DB,
		Scanner:         &resumeScanner{},
		SelectManyQuery: db.CreateSelectManyQuery(tableName, scanFields),
		SelectQuery:     db.CreateSelectQuery(tableName, scanFields, whereClause),
		CreateQuery:     db.CreateCreateQuery(tableName, neededFields, scanFields),
		UpdateQuery:     db.CreateUpdateQuery(tableName, updateFields, scanFields, whereClause),
		DeleteQuery:     db.CreateDeleteQuery(tableName),
	}
}

type resumeScanner struct{}

func (s *resumeScanner) Scan(row db.Scannable) (types.Resume, error) {
	var r types.Resume
	return r, row.Scan(
		&r.Id,
		&r.UserId,
		&r.Name,
		&r.Note,
		&r.CreatedAt,
		&r.UpdatedAt,
	)
}
