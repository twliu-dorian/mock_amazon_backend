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

	created, err = Get(obj.UserId)
	if err != nil {
		err = new(apierror.ApiError).From(err)
		return
	}

	return
}

func Get(userId string) (obj *User, err error) {
	obj, err = dao.Get(userId)
	if err != nil {
		err = new(apierror.ApiError).From(err)
		return
	}

	return
}
