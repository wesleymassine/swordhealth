package domain

import "errors"

var (
	ErrNoRowsResult = errors.New("sql: no rows in result set")
	ErrNotFound     = errors.New("Not found")
)
