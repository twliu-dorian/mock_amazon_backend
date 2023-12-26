package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"mock_amazon_backend/apierror"

	"github.com/asaskevich/govalidator"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/labstack/echo/v4"
)

func SuccessResponse(c echo.Context) (err error) {
	res := map[string]bool{"success": true}
	err = c.JSON(http.StatusOK, res)
	if err != nil {
		err = new(apierror.ApiError).From(err)
		return
	}

	data, _ := json.Marshal(res)
	c.Set("response", string(data))

	return
}

func JSONResponse(c echo.Context, response interface{}) (err error) {
	err = c.JSON(http.StatusOK, response)
	if err != nil {
		err = new(apierror.ApiError).From(err)
		return
	}

	data, _ := json.Marshal(response)
	c.Set("response", string(data))

	return
}

func ListResponseFromNoPaging(c echo.Context, page, limit int, input []interface{}) (err error) {
	length := len(input)
	c.Response().Header().Set("Pagination-Page", strconv.Itoa(page))
	c.Response().Header().Set("Pagination-Limit", strconv.Itoa(limit))
	c.Response().Header().Set("Pagination-Count", strconv.Itoa(length))

	var objects []interface{}
	if page*limit >= length {
		objects = make([]interface{}, 0)
	} else {
		end := (page + 1) * limit
		if end > length {
			end = length
		}
		objects = input[page*limit : end]
	}

	response := struct {
		Data []interface{} `json:"data"`
	}{
		Data: objects,
	}

	err = c.JSON(http.StatusOK, response)
	if err != nil {
		err = new(apierror.ApiError).From(err)
		return
	}

	data, _ := json.Marshal(response)
	c.Set("response", string(data))

	return
}

func ListResponse(c echo.Context, total int64, page, limit int, objects interface{}) (err error) {
	c.Response().Header().Set("Pagination-Page", strconv.Itoa(page))
	c.Response().Header().Set("Pagination-Limit", strconv.Itoa(limit))
	c.Response().Header().Set("Pagination-Count", strconv.FormatInt(total, 10))

	response := struct {
		Data interface{} `json:"data"`
	}{
		Data: objects,
	}

	err = c.JSON(http.StatusOK, response)
	if err != nil {
		err = new(apierror.ApiError).From(err)
		return
	}

	data, _ := json.Marshal(response)
	c.Set("response", string(data))

	return
}

func ListInRuleParams(inputs ...interface{}) validation.InRule {
	result := make([]interface{}, 0)
	for _, input := range inputs {
		i, isNil := validation.Indirect(input)
		if !isNil {
			result = append(result, i)
		}
	}

	return validation.In(result...)
}

func HTTPJSONRequest(method, url string, requestBody interface{}, headers map[string]string) (statusCode int, response []byte, err error) {
	var buffer io.Reader
	if headers == nil {
		headers = make(map[string]string)
	}

	if requestBody != nil {
		headers["Content-Type"] = "application/json"

		buffer = new(bytes.Buffer)
		err = json.NewEncoder(buffer.(*bytes.Buffer)).Encode(requestBody)
		if err != nil {
			err = new(apierror.ApiError).From(err)
			return
		}
	}

	return HTTPRequest(method, url, buffer, headers)
}

func HTTPRequest(method, url string, buffer io.Reader, headers map[string]string) (statusCode int, response []byte, err error) {
	if err = validation.Validate(&method, ListInRuleParams(
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	)); err != nil {
		err = new(apierror.ApiError).From(err)
		return
	}

	if err = validation.Validate(&url, validation.Required, is.URL); err != nil {
		err = new(apierror.ApiError).From(err)
		return
	}

	req, err := http.NewRequest(method, url, buffer)
	if err != nil {
		err = new(apierror.ApiError).From(err)
		return
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := new(http.Client).Do(req)
	if err != nil {
		err = new(apierror.ApiError).From(err)
		return
	}
	defer resp.Body.Close()

	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = new(apierror.ApiError).From(err)
		return
	}

	statusCode = resp.StatusCode

	if resp.StatusCode != 200 {
		err = new(apierror.ApiError).From(fmt.Errorf("[%s] %s: %d, %s", method, url, resp.StatusCode, string(response)))
		return
	}

	return
}

var currentTimeFunc func() time.Time

func SetCurrentTimeFunc(f func() time.Time) {
	currentTimeFunc = f
}

func Now() time.Time {
	if currentTimeFunc == nil {
		currentTimeFunc = time.Now
	}

	return currentTimeFunc()
}

var IsHex validation.StringRule

func init() {
	IsHex = validation.NewStringRuleWithError(govalidator.IsHexadecimal, validation.NewError("validation_is_hexadecimal", "must be a valid hexadecimal string"))
}
