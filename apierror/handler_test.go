package apierror

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	num := 0
	err := validation.Validate(&num, validation.Required)
	require.Error(t, err)

	echoError := errorHandler(err)
	fmt.Println(echoError)

	err = new(ApiError).From(err)
	echoError = errorHandler(err)
	fmt.Println(echoError)

	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "https://www.example.com", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = validation.Validate(&num, validation.Required)
	HTTPErrorHandler(err, c)

	fmt.Println(rec.Code)
	response, err := ioutil.ReadAll(rec.Body)
	require.NoError(t, err)
	fmt.Println(string(response))
}
