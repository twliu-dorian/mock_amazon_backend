package user

type Request struct {
	UserId string    `json:"adminId"`
	Email   string    `json:"email"`
	Salt         util.Bytes    `json:"-" db:"salt"`
	PasswordHash util.Bytes    `json:"-" db:"password_hash"`
	Role    AdminRole `json:"role"`
}