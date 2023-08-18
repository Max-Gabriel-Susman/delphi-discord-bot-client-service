package delphierror

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
)

type statusCoder struct {
	error
	code int
}

func WithStatusCode(err error, code int) error {
	if err == nil {
		err = errors.New(http.StatusText(code))
	}
	code = checkShouldOverride(err, code)
	return statusCoder{err, code}
}

func checkShouldOverride(err error, code int) int {
	var timeouter interface {
		error
		Timeout() bool
	}
	if errors.As(err, &timeouter) && timeouter.Timeout() {
		return http.StatusGatewayTimeout
	}
	var temper interface {
		error
		Temporary() bool
	}
	if errors.As(err, &temper) && temper.Temporary() {
		return http.StatusServiceUnavailable
	}
	if errors.Is(err, context.Canceled) {
		return 499
	}
	if errors.Is(err, sql.ErrNoRows) {
		return http.StatusNotFound
	}
	return code
}
