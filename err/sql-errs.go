package errs

import (
	"errors"
	"net/http"

	"github.com/jackc/pgconn"
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
	return http.StatusInternalServerError
}
