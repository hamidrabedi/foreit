package models

import (
	"fmt"
	"strings"
)

// QueryExpr represents a complex query with AND/OR/NOT logic.
// Use And, Or, and Not to build complex queries.
type QueryExpr struct {
	And       []QueryExpr
	Or        []QueryExpr
	Not       bool
	Condition Condition
}

// ToSQL converts the QueryExpr to SQL.
func (q *QueryExpr) ToSQL() (string, []interface{}) {
	if len(q.And) > 0 {
		return q.buildAnd()
	}
	
	if len(q.Or) > 0 {
		return q.buildOr()
	}
	
	if q.Condition != nil {
		sql, args := q.Condition.ToSQL()
		if q.Not {
			return fmt.Sprintf("NOT (%s)", sql), args
		}
		return sql, args
	}
	
	return "1=1", nil
}

func (q *QueryExpr) buildAnd() (string, []interface{}) {
	var parts []string
	var allArgs []interface{}
	
	for _, subQ := range q.And {
		sql, args := subQ.ToSQL()
		parts = append(parts, fmt.Sprintf("(%s)", sql))
		allArgs = append(allArgs, args...)
	}
	
	combinedSQL := strings.Join(parts, " AND ")
	if q.Not {
		combinedSQL = fmt.Sprintf("NOT (%s)", combinedSQL)
	}
	
	return combinedSQL, allArgs
}

func (q *QueryExpr) buildOr() (string, []interface{}) {
	var parts []string
	var allArgs []interface{}
	
	for _, subQ := range q.Or {
		sql, args := subQ.ToSQL()
		parts = append(parts, fmt.Sprintf("(%s)", sql))
		allArgs = append(allArgs, args...)
	}
	
	combinedSQL := strings.Join(parts, " OR ")
	if q.Not {
		combinedSQL = fmt.Sprintf("NOT (%s)", combinedSQL)
	}
	
	return combinedSQL, allArgs
}

// And creates an AND condition from multiple QueryExpr queries.
func And(conditions ...QueryExpr) QueryExpr {
	return QueryExpr{And: conditions}
}

// Or creates an OR condition from multiple QueryExpr queries.
func Or(conditions ...QueryExpr) QueryExpr {
	return QueryExpr{Or: conditions}
}

// Not creates a NOT condition from a QueryExpr.
func Not(condition QueryExpr) QueryExpr {
	return QueryExpr{Not: true, Condition: condition.Condition}
}

