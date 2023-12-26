package apierror

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"mock_amazon_backend/log"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func errorHandler(err error) (echoError *echo.HTTPError) {
	var apiError *ApiError
	var ok bool
	if apiError, ok = err.(*ApiError); ok {
		err = apiError.err
	}

	if apiError == nil {
		apiError = new(ApiError).From(err, 3)
	}

	// ozzo-validation errors
	if _, ok = err.(validation.ErrorObject); ok {
		echoError = echo.NewHTTPError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest)).SetInternal(apiError)
	}
	if _, ok = err.(validation.Errors); ok {
		echoError = echo.NewHTTPError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest)).SetInternal(apiError)
	}

	// SQL errors
	if errors.Is(err, sql.ErrNoRows) {
		echoError = echo.NewHTTPError(http.StatusNotFound, http.StatusText(http.StatusNotFound)).SetInternal(apiError)
	}
	var mse *mysql.MySQLError
	if errors.As(err, &mse) && mse.Number == 1062 {
		echoError = echo.NewHTTPError(http.StatusConflict, http.StatusText(http.StatusConflict)).SetInternal(apiError)
	}

	if errors.As(err, &echoError) {
		return
	}

	if echoError == nil {
		switch apiError.Code {
		case FORBIDDEN:
			echoError = &echo.HTTPError{
				Code:     http.StatusForbidden,
				Message:  http.StatusText(http.StatusForbidden),
				Internal: apiError,
			}
		case AUTHENTICATION:
			echoError = &echo.HTTPError{
				Code:     http.StatusUnauthorized,
				Message:  http.StatusText(http.StatusUnauthorized),
				Internal: apiError,
			}
		case INPUT_ERROR:
			echoError = &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  http.StatusText(http.StatusBadRequest),
				Internal: apiError,
			}
		case CONFLICT:
			echoError = &echo.HTTPError{
				Code:     http.StatusConflict,
				Message:  http.StatusText(http.StatusConflict),
				Internal: apiError,
			}
		case UPSTREAM_ERROR:
			echoError = &echo.HTTPError{
				Code:     http.StatusServiceUnavailable,
				Message:  http.StatusText(http.StatusServiceUnavailable),
				Internal: apiError,
			}
		case NOT_FOUND:
			echoError = &echo.HTTPError{
				Code:     http.StatusNotFound,
				Message:  http.StatusText(http.StatusNotFound),
				Internal: apiError,
			}
		default:
			echoError = &echo.HTTPError{
				Code:     http.StatusInternalServerError,
				Message:  http.StatusText(http.StatusInternalServerError),
				Internal: apiError,
			}
		}
	}

	return
}

// Adapted from https://pkg.go.dev/github.com/labstack/echo/v4#Echo.DefaultHTTPErrorHandler
func HTTPErrorHandler(err error, c echo.Context) {
	he := errorHandler(err)

	// Issue #1426
	code := he.Code
	message := he.Message
	if m, ok := message.(string); ok {
		message = echo.Map{
			"error": echo.Map{"message": m},
		}
	}

	logLabel := log.LabelMonitor
	switch apiError := he.Internal.(type) {
	case *ApiError:
		if apiError.Label != "" {
			logLabel = apiError.Label
		}
	}

	data, _ := json.Marshal(message)
	log.Error(
		logLabel,
		fmt.Sprintf("Status=%d", code),
		" ClientAddr=", c.RealIP(),
		" Response=", string(data),
		" Internal=", he.Internal,
	)

	// Send response
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead { // Issue #608
			err = c.NoContent(he.Code)
		} else {
			err = c.JSON(code, message)
		}
		if err != nil {
			log.Error(log.LabelMonitor, err)
		}
	}
}
