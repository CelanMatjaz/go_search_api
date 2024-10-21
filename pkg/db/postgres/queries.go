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
	createWithTags    func(tagCount int) string
}

func createQueryHolder[T any](table string, hasPagination bool) QueryHolder {
	return QueryHolder{
		querySingle:  SingleRecordQuery[T](table),
		queryMany:    ManyRecordsQuery[T](table, hasPagination),
		createSingle: CreateRecordQuery[T](table),
		updateSingle: UpdateRecordQuery[T](table),
		deleteSingle: DeleteRecordQuery(table),
	}
}

func createQueryHolderWithTags[T any](table string, mtmTable string) QueryHolderWithCreateTags {
	return QueryHolderWithCreateTags{
		QueryHolder:       createQueryHolder[T](table, true),
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
		return strings.Join(indices, ", ")
	},
	"joinWithPrepend": func(fields []string, prepend string) string {
		values := make([]string, len(fields))
		for i, val := range fields {
			values[i] = fmt.Sprintf("%s%s", prepend, val)
		}
		return strings.Join(values, ", ")
	},
	"setFields": func(fields []string, start int) string {
		values := make([]string, 0)
		skip := 0
		for i, field := range fields {
			switch field {
			case "id":
				fallthrough
			case "account_id":
				fallthrough
			case "created_at":
				skip += 1
				break
			case "updated_at":
				values = append(values, fmt.Sprintf("%s = DEFAULT", field))
				skip += 1
				break
			default:
				values = append(values, fmt.Sprintf("%s = $%d", field, i+start-skip))
			}
		}

		return strings.Join(values, ", ")
	},
	"lenMoreThan0": func(value []string) bool { return len(value) > 0 },
}

type BasicQueryData struct {
	RecordTable string
	Fields      []string
}

func execQueryTemplate(templateStr string, data any) string {
	t := template.Must(template.New("").Funcs(queryFuncMap).Parse(templateStr))

	var result bytes.Buffer
	err := t.Execute(&result, data)
	if err != nil {
		panic(fmt.Sprintf("could not generate query from template\n\terror: %s", err.Error()))
	}

	return result.String()
}

func SingleRecordQuery[T any](recordTable string) string {
	tmpl := `
        SELECT {{ join .Fields ", " }} 
        FROM {{ .RecordTable }} 
        WHERE account_id = $1 AND id = $2`

	return execQueryTemplate(tmpl, BasicQueryData{
		RecordTable: recordTable,
		Fields:      GetDbFields[T](),
	})
}

func ManyRecordsQuery[T any](recordTable string, addPagination bool) string {
	tmpl := `
        SELECT {{ join .Fields ", " }} 
        FROM {{ .RecordTable }} 
        WHERE account_id = $1
        {{ if .Pagination }} OFFSET $2 LIMIT $3 {{ end }}`

	type Data struct {
		RecordTable string
		Fields      []string
		Pagination  bool
	}

	return execQueryTemplate(tmpl, Data{
		RecordTable: recordTable,
		Fields:      GetDbFields[T](),
		Pagination:  addPagination,
	})
}

func CreateRecordQuery[T any](recordTable string) string {
	tmpl := `
        INSERT INTO {{ .RecordTable }} ({{ join .Fields ", " }})
        VALUES ({{ joinIndices .Fields 1 "" }})
        RETURNING *`

	return execQueryTemplate(tmpl, BasicQueryData{
		RecordTable: recordTable,
		Fields:      GetDbFieldsForCreate[T](),
	})
}

func UpdateRecordQuery[T any](recordTable string) string {
	tmpl := `
        UPDATE {{ .RecordTable }}
        SET {{ setFields .Fields 3 }}
        WHERE id = $1 AND account_id = $2
        RETURNING *`

	return execQueryTemplate(tmpl, BasicQueryData{
		RecordTable: recordTable,
		Fields:      GetDbFieldsForUpdate[T](),
	})
}

func DeleteRecordQuery(recordTable string) string {
	tmpl := `
        DELETE FROM {{ .RecordTable }}
        WHERE id = $1 AND account_id = $2`

	type DeleteQueryData struct {
		RecordTable string
	}

	return execQueryTemplate(tmpl, DeleteQueryData{
		RecordTable: recordTable,
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
        ORDER BY lr.updated_at DESC`

	type Data struct {
		BasicQueryData
		TagFields []string
		MtmTable  string
	}

	return execQueryTemplate(tmpl, Data{
		BasicQueryData: BasicQueryData{
			RecordTable: recordTable,
			Fields:      GetDbFields[T](),
		},
		TagFields: GetDbFields[types.Tag](),
		MtmTable:  mtmTable,
	})
}

func CreateCreateSingleWithTags[T any](recordTable string, mtmTable string) func(tagCount int) string {
	baseTmpl := `
        INSERT INTO {{ .RecordTable }} ({{ join .Fields ", " }}) 
        VALUES ({{ joinIndices .Fields 1 "" }})
        RETURNING *`

	fields := GetDbFieldsForCreate[T]()
	baseInsertQuery := execQueryTemplate(baseTmpl, BasicQueryData{
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

	leftSide := execQueryTemplate(leftTmpl, Data{
		BaseQuery: baseInsertQuery,
		MtmTable:  mtmTable,
	})

	rightSide := `) SELECT cs.* FROM new_record cs`

	start := len(fields)
	generateTagInsertValues := func(tagCount int, start int) string {
		sections := make([]string, tagCount)
		for i := 0; i < tagCount; i++ {
			sections[i] = fmt.Sprintf("($%d, (SELECT id FROM new_record))", i+start+1)
		}
		return strings.Join(sections, ", ")
	}

	return func(tagCount int) string {
		if tagCount == 0 {
			return baseInsertQuery
		}

		return fmt.Sprintf("%s %s %s",
			leftSide,
			generateTagInsertValues(tagCount, start),
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

	return execQueryTemplate(tmpl, Data{
		MtmTable: mtmTable,
		Fields:   GetDbFields[types.Tag](),
	})
}
