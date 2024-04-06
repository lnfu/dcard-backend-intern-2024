package utils

import (
	"database/sql"
)

// TODO test
func NullInt32FromInt32Pointer(int32_p *int32) sql.NullInt32 {
	if int32_p == nil {
		return sql.NullInt32{Int32: 0, Valid: false}
	}
	return sql.NullInt32{Int32: int32(*int32_p), Valid: true}
}
