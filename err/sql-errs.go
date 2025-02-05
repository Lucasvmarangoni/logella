package errs

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
)

func GetHTTPStatusFromPgError(err error) int {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code[:2] {
		case "23":
			return http.StatusBadRequest
		case "42":
			return http.StatusInternalServerError
		case "53":
			return http.StatusServiceUnavailable
		default:
			return http.StatusInternalServerError
		}
	}
	if err.Error() == "no rows in result set" {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}
