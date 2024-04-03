package utils

import (
	"database/sql"
	"strconv"
	"strings"
)

func NonEmptyNullStringFromString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func NullInt32FromString(s string) (sql.NullInt32, error) {
	if strings.TrimSpace(s) == "" {
		return sql.NullInt32{Valid: false}, nil
	}
	temp, err := strconv.Atoi(s)
	if err != nil {
		return sql.NullInt32{Valid: false}, err
	}
	return sql.NullInt32{Int32: int32(temp), Valid: true}, nil
}

// TODO test
func NullInt32FromInt32Pointer(int_p *int) sql.NullInt32 {
	if int_p == nil {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: int32(*int_p), Valid: true}
}
