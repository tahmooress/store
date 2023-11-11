package psql

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidParentheses = errors.New("invalid parentheses")

type queryBuilder struct {
	baseQuery      string
	args           []interface{}
	builder        *strings.Builder
	idx            int
	whereClzPlaced bool
	prStack        stack
	err            error
}

func newQueryBuilder(baseQuery string) *queryBuilder {
	qb := &queryBuilder{
		baseQuery: baseQuery,
		args:      make([]interface{}, 0),
		builder:   new(strings.Builder),
		idx:       1,
	}
	_, err := qb.builder.WriteString(baseQuery)
	if err != nil {
		qb.err = err
	}

	return qb
}

func (q *queryBuilder) reset() {
	q.baseQuery = ""
	q.args = nil
	q.builder.Reset()
	q.idx = 1
	q.prStack = nil
	q.whereClzPlaced = false
	q.err = nil
}

func (q *queryBuilder) where() *queryBuilder {
	if q.whereClzPlaced {
		return q
	}

	_, err := q.builder.WriteString(" WHERE")
	if err != nil {
		q.err = err
	}
	q.whereClzPlaced = true
	return q
}

func (q *queryBuilder) stratParentheses() *queryBuilder {
	if _, err := q.builder.WriteString(" ("); err != nil {
		q.err = err
	}

	q.prStack.push("(")
	return q
}

func (q *queryBuilder) closeParentheses() *queryBuilder {
	if _, err := q.builder.WriteString(" )"); err != nil {
		q.err = err
	}

	q.prStack.push(")")
	return q
}

func (q *queryBuilder) whereCond(key string, val interface{}) *queryBuilder {
	q.where()
	if _, err := q.builder.WriteString(fmt.Sprintf(" %s", key)); err != nil {
		q.err = err
	}
	if _, err := q.builder.WriteString(" ="); err != nil {
		q.err = err
	}
	if _, err := q.builder.WriteString(q.nextPlaceHolder()); err != nil {
		q.err = err
	}
	q.args = append(q.args, val)
	q.whereClzPlaced = true
	return q
}

func (q *queryBuilder) wherein(key string, vals []interface{}) *queryBuilder {
	if len(vals) == 0 {
		return q
	}
	q.where()
	if _, err := q.builder.WriteString(fmt.Sprintf(" %s", key)); err != nil {
		q.err = err
	}
	if _, err := q.builder.WriteString(" IN ("); err != nil {
		q.err = err
	}

	var placeHolders []string
	for _, v := range vals {
		placeHolders = append(placeHolders, q.nextPlaceHolder())
		q.args = append(q.args, v)
	}

	if _, err := q.builder.WriteString(strings.Join(placeHolders, ",")); err != nil {
		q.err = err
	}

	if _, err := q.builder.WriteString(")"); err != nil {
		q.err = err
	}

	return q
}

func (q *queryBuilder) and() *queryBuilder {
	if _, err := q.builder.WriteString(" AND"); err != nil {
		q.err = err
	}
	return q
}

func (q *queryBuilder) or() *queryBuilder {
	if _, err := q.builder.WriteString(" OR"); err != nil {
		q.err = err
	}
	return q
}

func (q *queryBuilder) nextPlaceHolder() string {
	r := fmt.Sprintf(" $%d", q.idx)
	q.idx++
	return r
}

func (q *queryBuilder) appendQuery(str string) *queryBuilder {
	_, err := q.builder.WriteString(fmt.Sprintf(" %s", strings.TrimPrefix(str, " ")))
	if err != nil {
		q.err = err
	}
	return q
}

func (q *queryBuilder) build() (string, []interface{}, error) {
	if err := parseParentheses(q.prStack); err != nil {
		return "", nil, fmt.Errorf("queryBuilder build: %s", err)
	}

	return q.builder.String(), q.args, q.err
}

func parseParentheses(stk []string) error {
	var src stack

	for _, v := range stk {
		if src.isEmpty() {
			src.push(v)
		} else {
			t, _ := src.pop()
			if !(t == "(" && v == ")") {
				src.push(t)
				src.push(v)
			}
		}
	}

	if !src.isEmpty() {
		return ErrInvalidParentheses
	}

	return nil
}

type stack []string

func (s *stack) push(val string) {
	*s = append(*s, val)
}

func (s *stack) isEmpty() bool {
	return len(*s) == 0
}

func (s *stack) pop() (string, bool) {
	if s.isEmpty() {
		return "", false
	}

	index := len(*s) - 1
	element := (*s)[index]
	*s = (*s)[:index]
	return element, true
}
