package user

import (
	"mock_amazon_backend/apierror"
	"mock_amazon_backend/util"

	"github.com/google/uuid"
)

var dao daoInterface = new(daoImplement)

func Create(obj *User) (created *User, err error) {
	obj.UserId = uuid.New().String()
	obj.CreatedAt.SetTime(util.Now())
	obj.UpdatedAt.SetTime(util.Now())

	err = dao.Create(obj)

	if err != nil {
		err = new(apierror.ApiError).From(err)
		return
	}

	created, err = Get(obj.Email)
	if err != nil {
		err = new(apierror.ApiError).From(err)
		return
	}

	return
}

func Get(email string) (obj *User, err error) {
	obj, err = dao.Get(email)
	if err != nil {
		err = new(apierror.ApiError).From(err)
		return
	}

	return
}

func Login(request LoginRequest) (response LoginResponse, err error) {
	user, err := Get(request.Email)
	if err != nil {
		err = new(apierror.ApiError).From(err)
		return
	}
	if !user.ValidatePassword(request.Password) {
		err = new(apierror.ApiError).FromMessage("invalid password").SetCode(apierror.AUTHENTICATION)
		response.Valid = false
		return
	}
	response.Valid = true

	return
}
