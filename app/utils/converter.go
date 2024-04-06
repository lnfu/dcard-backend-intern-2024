package utils

import (
	"database/sql"
)

// TODO test
func NullInt32FromInt32Pointer(int_p *int) sql.NullInt32 {
	if int_p == nil {
		return sql.NullInt32{Int32: 0, Valid: false}
	}
	return sql.NullInt32{Int32: int32(*int_p), Valid: true}
}
