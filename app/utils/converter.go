package utils

import (
	"database/sql"
)

func NullInt32FromInt32Pointer(int32_p *int32) sql.NullInt32 {
	if int32_p == nil {
		return sql.NullInt32{Int32: 0, Valid: false}
	}
	return sql.NullInt32{Int32: (*int32_p), Valid: true}
}

func NullStringFromStringPointer(string_p *string) sql.NullString {
	if string_p == nil {
		return sql.NullString{String: "", Valid: false}
	}
	return sql.NullString{String: (*string_p), Valid: true}
}
