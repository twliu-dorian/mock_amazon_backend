package user

import (
	"mock_amazon_backend/apierror"
	"mock_amazon_backend/util"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
)

type Controller struct{}

func (Controller) Create(c echo.Context) (err error) {
	var request struct {
		Request
		Password string `json:"password"`
	}

	if err = c.Bind(&request); err != nil {
		err = new(apierror.ApiError).From(err).SetCode(apierror.INPUT_ERROR)
		return
	}
	if err = request.Validate(); err != nil {
		err = new(apierror.ApiError).From(err).SetCode(apierror.INPUT_ERROR)
		return
	}
	if err = validation.Validate(&request.Password, validation.Required); err != nil {
		err = new(apierror.ApiError).From(err).SetCode(apierror.INPUT_ERROR)
		return
	}
	obj := request.ToUser()
	err = obj.SetupPassword(request.Password)
	if err != nil {
		err = new(apierror.ApiError).From(err)
		return
	}

	response, err := Create(obj)
	if err != nil {
		err = new(apierror.ApiError).From(err)
		return
	}

	return util.JSONResponse(c, response)
}

func (Controller) Login(c echo.Context) (err error) {
	var request LoginRequest

	if err = c.Bind(&request); err != nil {
		err = new(apierror.ApiError).From(err).SetCode(apierror.INPUT_ERROR)
		return
	}

	if err = validation.Validate(&request.Password, validation.Required); err != nil {
		err = new(apierror.ApiError).From(err).SetCode(apierror.INPUT_ERROR)
		return
	}
	response, err := Login(request)
	if err != nil {
		err = new(apierror.ApiError).From(err)
		return
	}

	return util.JSONResponse(c, response)
}
