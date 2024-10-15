package postgres

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

type QueryHolder struct {
	querySingle  string
	queryMany    string
	createSingle string
	updateSingle string
	deleteSingle string
}

type QueryHolderWithCreateTags struct {
	QueryHolder
	queryManyWithTags string
	createWithTags    func(tagIds []int) string
}

func createQueryHolder[T any, bodyT any](table string, hasPagination bool) QueryHolder {
	return QueryHolder{
		querySingle:  SingleRecordQuery[T](table),
		queryMany:    ManyRecordsQuery[T](table, hasPagination),
		createSingle: CreateRecordQuery[bodyT](table),
		updateSingle: UpdateRecordQuery[bodyT](table),
		deleteSingle: DeleteRecordQuery[T](table),
	}
}

func createQueryHolderWithTags[T any, bodyT any](table string, mtmTable string) QueryHolderWithCreateTags {
	return QueryHolderWithCreateTags{
		QueryHolder:       createQueryHolder[T, bodyT](table, true),
		queryManyWithTags: ManyRecordsWithTagsQuery[T](table, mtmTable),
		createWithTags:    CreateCreateSingleWithTags[T](table, mtmTable),
	}
}

var queryFuncMap = template.FuncMap{
	"join": strings.Join,
	"joinIndices": func(fields []string, start int, prepend string) string {
		indices := make([]string, len(fields))
		for i := range indices {
			indices[i] = fmt.Sprintf("%s$%d", prepend, i+start)
		}
		return strings.Join(indices, ",")
	},
	"joinWithPrepend": func(fields []string, prepend string) string {
		values := make([]string, len(fields))
		for i, val := range fields {
			values[i] = fmt.Sprintf("%s%s", prepend, val)
		}
		return strings.Join(values, ", ")
	},
	"setFields": func(fields []string, start int) string {
		indices := make([]string, len(fields))
		skip := 0
		for i, field := range indices {
			switch field {
			case "id":
			case "account_id":
			case "created_at":
				skip += 1
				break
			case "updated_at":
				indices[i] = fmt.Sprintf("%s = DEFAULT", field)
				skip += 1
				break
			default:
				indices[i] = fmt.Sprintf("%s = $%d", field, i+start-skip)
			}

		}
		return strings.Join(indices, ", ")
	},
}

type BasicQueryData struct {
	RecordTable string
	Fields      []string
}

func execTemplate(templateStr string, data any) string {
	t := template.Must(template.New("").Funcs(queryFuncMap).Parse(templateStr))

	var result bytes.Buffer
	err := t.Execute(&result, data)
	if err != nil {
		panic("could not generate query from template")
	}

	return result.String()
}

func SingleRecordQuery[T any](recordTable string) string {
	tmpl := `
        SELECT {{ join .Fields "," }} 
        FROM {{ .RecordTable }} 
        WHERE account_id = $1, id = $2`

	return execTemplate(tmpl, BasicQueryData{
		RecordTable: recordTable,
		Fields:      reflectDbFields[T](),
	})
}

func ManyRecordsQuery[T any](recordTable string, addPagination bool) string {
	tmpl := `
        SELECT {{ join .Fields "," }} 
        FROM {{ .RecordTable }} 
        WHERE account_id = $1
        {{ if .Pagination }} OFFSET $2 LIMIT $3 {{ end }}`

	type Data struct {
		RecordTable string
		Fields      []string
		Pagination  bool
	}

	return execTemplate(tmpl, Data{
		RecordTable: recordTable,
		Fields:      reflectDbFields[T](),
		Pagination:  addPagination,
	})
}

func CreateRecordQuery[T any](recordTable string) string {
	tmpl := `
        INSERT INTO {{ .RecordTable }} ({{ join .Fields "," }})
        VALUES ({{ joinIndices .Fields 1 "" }})
        RETURNING *`

	return execTemplate(tmpl, BasicQueryData{
		RecordTable: recordTable,
		Fields:      reflectDbFields[T](),
	})
}

func UpdateRecordQuery[T any](recordTable string) string {
	tmpl := `
        UPDATE {{ .RecordTable }}
        SET {{ setFields .Fields 1 }}
        VALUES ({{ joinIndices .Fields 3 "" }})
        WHERE account_id = $1, id = $2
        RETURNING *`

	return execTemplate(tmpl, BasicQueryData{
		RecordTable: recordTable,
		Fields:      reflectDbFields[T](),
	})
}

func DeleteRecordQuery[T any](recordTable string) string {
	tmpl := `
        DELETE {{ .RecordTable }}
        WHERE FROM account_id = $1, id = $2
        RETURNING *`

	return execTemplate(tmpl, BasicQueryData{
		RecordTable: recordTable,
		Fields:      reflectDbFields[T](),
	})
}

func ManyRecordsWithTagsQuery[T any](recordTable string, mtmTable string) string {
	tmpl := ` 
        WITH limited_records AS (
            SELECT {{ join .Fields ", " }}
            FROM {{ .RecordTable }}
            WHERE account_id = $1
            OFFSET $2 LIMIT $3
        )
        SELECT {{ joinWithPrepend .Fields "lr." }}, {{ joinWithPrepend .TagFields "t." }}
        FROM limited_records lr
        LEFT JOIN {{ .MtmTable }} ta ON lr.id = ta.record_id
        LEFT JOIN tags t ON ta.tag_id = t.id
        ORDER BY lr.updated_at DESC;`

	type Data struct {
		BasicQueryData
		TagFields []string
		MtmTable  string
	}

	return execTemplate(tmpl, Data{
		BasicQueryData: BasicQueryData{
			RecordTable: recordTable,
			Fields:      reflectDbFields[T](),
		},
		TagFields: reflectDbFields[types.Tag](),
		MtmTable:  mtmTable,
	})
}

func CreateCreateSingleWithTags[T any](recordTable string, mtmTable string) func(tagIds []int) string {
	baseTmpl := `
        INSERT INTO {{ .RecordTable }} ({{ join .Fields "," }}) 
        VALUES ({{ joinIndices .Fields 1 "" }})
        RETURNING *`

	fields := reflectDbFields[T]()
	baseInsertQuery := execTemplate(baseTmpl, BasicQueryData{
		RecordTable: recordTable,
		Fields:      fields,
	})

	leftTmpl := `
        WITH new_record AS (
            {{ .BaseQuery }}
        ),
        _ AS (
            INSERT INTO {{ .MtmTable }} (tag_id, record_id) 
            VALUES`

	type Data struct {
		BaseQuery string
		MtmTable  string
	}

	leftSide := execTemplate(leftTmpl, Data{
		BaseQuery: baseInsertQuery,
		MtmTable:  mtmTable,
	})

	rightSide := `) SELECT cs.* FROM new_record cs`

	start := len(fields)
	generateTagInsertValues := func(tagIds []int, start int) string {
		sections := make([]string, len(tagIds))
		for i := range tagIds {
			sections[i] = fmt.Sprintf("($%d, (SELECT id FROM new_record))", i+start+1)
		}
		return strings.Join(sections, ",")
	}

	return func(tagIds []int) string {
		if len(tagIds) == 0 {
			return baseInsertQuery
		}

		return fmt.Sprintf("%s %s %s",
			leftSide,
			generateTagInsertValues(tagIds, start),
			rightSide,
		)
	}
}

func createTagJoinQuery(mtmTable string) string {
	tmpl := `
        SELECT {{ joinIndices .Fields 3 "" }} 
        FROM tags 
        LEFT JOIN {{ .MtmTable }} t ON t.tag_id = tag.id 
        WHERE tags.account_id = $1 AND t.record_id = $2`

	type Data struct {
		MtmTable string
		Fields   []string
	}

	return execTemplate(tmpl, Data{
		MtmTable: mtmTable,
		Fields:   reflectDbFields[types.Tag](),
	})
}
