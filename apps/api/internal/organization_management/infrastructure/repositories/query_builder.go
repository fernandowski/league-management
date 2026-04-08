package repositories

import (
	"fmt"
	"strings"
)

type QueryBuilder struct {
	baseQuery    string
	whereClauses []string
	parameters   []interface{}
	limit        int
	offset       int
}

func (qb *QueryBuilder) SetQuery(query string) {
	qb.baseQuery = query
}

func (qb *QueryBuilder) AddFilter(condition string, param interface{}) {
	qb.whereClauses = append(qb.whereClauses, condition)
	qb.parameters = append(qb.parameters, param)
}

func (qb *QueryBuilder) SetPagination(limit, offset int) {
	qb.limit = limit
	qb.offset = offset
}

func (qb *QueryBuilder) BuildQuery() (string, []interface{}) {
	query := qb.baseQuery

	if len(qb.whereClauses) > 0 {
		query += " WHERE " + strings.Join(qb.whereClauses, " AND ")
	}

	if qb.limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", qb.limit)
	}
	if qb.offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", qb.offset)
	}

	return query, qb.parameters
}
