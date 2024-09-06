package db

import (
	"fmt"
	"strconv"
	"strings"
)

func CreateSelectManyQuery(table string, fields []string) string {
	return fmt.Sprintf(
		`SELECT %s FROM %s WHERE user_id = $1 ORDER BY id DESC OFFSET $2 LIMIT $3`,
		strings.Join(fields, ", "), table,
	)
}

func CreateSelectQuery(table string, fields []string, whereClause string) string {
	return fmt.Sprintf(`
        SELECT %s FROM %s %s`,
		strings.Join(fields, ", "), table, whereClause,
	)
}

func CreateCreateQuery(table string, insertFields []string, allFields []string) string {
	return fmt.Sprintf(`
        INSERT INTO %s (%s)
        VALUES (%s) RETURNING %s`,
		table,
		strings.Join(insertFields, ", "),
		strings.Join(getArgs(insertFields), ", "),
		strings.Join(allFields, ", "),
	)
}

func CreateUpdateQuery(table string, fields []string, allFields []string) string {
	return fmt.Sprintf(`
        UPDATE %s
        SET %s
        WHERE id = $1 RETURNING %s`,
		table,
		strings.Join(getSetArgs(2, fields), ", "),
		strings.Join(allFields, ", "),
	)
}

func CreateDeleteQuery(table string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, table)
}

func getArgs(fields []string) []string {
	args := make([]string, 0)
	for i := 0; i < len(fields); i++ {
		args = append(args, "$"+strconv.Itoa(i+1))
	}
	return args
}

func getSetArgs(start int, fields []string) []string {
	args := make([]string, 0)
	for i := 0; i < len(fields); i++ {
		args = append(args, fmt.Sprintf("%s = $%d", fields[i], i+start))
	}
	return args
}
