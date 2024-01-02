package user

import (
	"mock_amazon_backend/apierror"
	"mock_amazon_backend/database"
)

type daoInterface interface {
	Create(user *User) error
	List(cond *Condition) (int64, []User, error)
	Get(id string) (*User, error)
	Update(toUpdate *User) error
	Delete(id string) error
}

type daoImplement struct{}

func (dao *daoImplement) Create(user *User) (err error) {
	db, err := database.DB()
	if err != nil {
		err = new(apierror.ApiError).From(err)
		return
	}

	tx, err := db.Beginx()
	if err != nil {
		err = new(apierror.ApiError).From(err)
		return
	}

	defer database.ClearTransition(tx)
	query := `
		INSERT INTO
			user (
				user_id,
				email,
				salt,
				password_hash,
				created_at,
				updated_at
			)
		VALUES
			(
				:user_id,
				:email,
				:salt,
				:password_hash,
				:created_at,
				:updated_at
			)`

	_, err = tx.NamedExec(query, user)
	if err != nil {
		err = new(apierror.ApiError).From(err)
		return
	}

	err = tx.Commit()
	if err != nil {
		err = new(apierror.ApiError).From(err)
		return
	}
	return
}

func (dao *daoImplement) List(cond *Condition) (total int64, users []User, err error) {

	return
}

func (dao *daoImplement) Get(id string) (user *User, err error) {

	return
}

func (dao *daoImplement) Update(toUpdate *User) (err error) {

	return
}

func (dao *daoImplement) Delete(id string) (err error) {

	return
}
