package apierror

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"
)

type ErrorCode int

const (
	INTERNAL_ERROR ErrorCode = iota
	UPSTREAM_ERROR
	INPUT_ERROR
	NOT_FOUND
	FORBIDDEN
	AUTHENTICATION
	CONFLICT
)

type LogLabel int

const (
	MONITOR LogLabel = iota
	SIGN
)

type ApiError struct {
	err     error
	Code    ErrorCode `json:"code"`
	Message string    `json:"message,omitempty"`
	Label   string    `json:"-"`
	Stacks  []string  `json:"stacks"`
}

func (err ApiError) MarshalJSON() (data []byte, e error) {
	var message string
	if err.err != nil {
		message = err.err.Error()
	}

	type apiError ApiError
	v := struct {
		*apiError
		Message string `json:"message"`
	}{
		apiError: (*apiError)(&err),
		Message:  message,
	}

	return json.Marshal(v)
}

func (err *ApiError) FromMessage(m string) *ApiError {
	return err.From(errors.New(m), 2)
}

func (err *ApiError) FromSprintf(format string, args ...interface{}) *ApiError {
	return err.From(fmt.Errorf(format, args...), 2)
}

func (err *ApiError) From(e error, skips ...int) *ApiError {
	skip := 1
	if len(skips) == 1 {
		skip = skips[0]
	}
	pc, _, _, _ := runtime.Caller(skip)
	f := runtime.FuncForPC(pc)
	file, line := f.FileLine(pc)
	stack := fmt.Sprintf("%s:%d", file, line)

	var apiError *ApiError
	if errors.As(e, &apiError) {
		err.Code = apiError.Code
		err.Message = apiError.Message
		err.err = apiError.err
		err.Stacks = append(apiError.Stacks, stack)
		return err
	}

	err.Code = INTERNAL_ERROR
	err.err = e
	err.Stacks = append(err.Stacks, stack)
	return err
}

func (err *ApiError) SetCode(code ErrorCode) *ApiError {
	err.Code = code
	return err
}

func (err *ApiError) SetMessage(msg string) *ApiError {
	err.Message = msg
	return err
}

func (err *ApiError) SetLabel(label string) *ApiError {
	err.Label = label
	return err
}

func (err ApiError) Error() string {
	stacks := strings.Join(err.Stacks, " ")
	return fmt.Sprintf("%s (%s)", err.err.Error(), stacks)
}

func (err ApiError) Unwrap() error {
	return err.err
}
