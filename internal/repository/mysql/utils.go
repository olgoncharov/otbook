package mysql

import "fmt"

func queryWithLimitOffset(query string, limit, offset uint) string {
	if limit == 0 {
		return query
	}

	if offset == 0 {
		return query + fmt.Sprintf(" LIMIT %d", limit)
	}

	return query + fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
}
