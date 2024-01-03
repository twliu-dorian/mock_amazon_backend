package user

import (
	"bytes"
	"crypto/rand"
	"crypto/sha512"

	"mock_amazon_backend/apierror"
	"mock_amazon_backend/util"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"golang.org/x/crypto/pbkdf2"
)

type Request struct {
	Email string `json:"email"`
}

type User struct {
	UserId       string         `json:"userId" db:"user_id"`
	Email        string         `json:"email" db:"email"`
	Salt         util.Bytes     `json:"-" db:"salt"`
	PasswordHash util.Bytes     `json:"-" db:"password_hash"`
	CreatedAt    util.Timestamp `json:"-" db:"created_at"`
	UpdatedAt    util.Timestamp `json:"-" db:"updated_at"`
}

func (r Request) Validate() (err error) {
	err = validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, is.EmailFormat),
	)
	return
}

func (r *Request) ToUser() *User {
	return &User{
		Email: r.Email,
	}
}

func (u *User) SetupPassword(password string) (err error) {
	u.Salt = make([]byte, 16)
	length, err := rand.Read(u.Salt)
	if err != nil {
		return
	}
	if length != 16 {
		err = new(apierror.ApiError).FromMessage("random generator error")
		return
	}

	u.PasswordHash = pbkdf2.Key([]byte(password), u.Salt, 10000, 64, sha512.New)

	return
}

func (u User) ValidatePassword(password string) bool {
	hash := pbkdf2.Key([]byte(password), u.Salt, 10000, 64, sha512.New)

	return bytes.Equal(u.PasswordHash, hash)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Valid bool `json:"valid"`
}

// paging condition
type Condition struct {
	util.PagingCondition
}
