package apierror

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	e := f()
	fmt.Println(e)
	data, err := json.MarshalIndent(e, "", "  ")
	require.NoError(t, err)
	fmt.Println(string(data))

	e1 := f1()
	fmt.Println(e1)
	data, err = json.MarshalIndent(e1, "", "  ")
	require.NoError(t, err)
	fmt.Println(string(data))
}

func f() error {
	return new(ApiError).FromMessage("test")
}

func f1() error {
	err := f()
	return new(ApiError).From(err)
}
