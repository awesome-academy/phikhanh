package utils

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

// PostgreSQL error codes
const (
	pgErrUniqueViolation     = "23505"
	pgErrForeignKeyViolation = "23503"
	pgErrNotNullViolation    = "23502"
)

// ParseDBError - Chuyển đổi DB error thành ServiceError với message rõ ràng
func ParseDBError(err error) *ServiceError {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return NewInternalServerError(err)
	}

	switch pgErr.Code {
	case pgErrUniqueViolation:
		field := extractFieldFromDetail(pgErr.Detail)
		return NewBadRequestError("The " + field + " already exists. Please use a different value.")

	case pgErrForeignKeyViolation:
		return NewBadRequestError("Related record not found. Please check your input.")

	case pgErrNotNullViolation:
		return NewBadRequestError("Field '" + pgErr.ColumnName + "' is required.")

	default:
		return NewInternalServerError(err)
	}
}

// extractFieldFromDetail - Lấy tên field từ pg error detail
// Detail format: `Key (code)=(SV001) already exists.`
func extractFieldFromDetail(detail string) string {
	start := strings.Index(detail, "(")
	end := strings.Index(detail, ")")
	if start == -1 || end == -1 || end <= start {
		return "value"
	}
	return detail[start+1 : end]
}
