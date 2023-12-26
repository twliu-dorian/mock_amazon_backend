package util

import (
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"time"

	"mock_amazon_backend/apierror"

	"github.com/shopspring/decimal"
)

type PagingCondition struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}

type Bytes []byte

func (b *Bytes) Scan(value interface{}) (err error) {
	var v sql.NullString
	err = v.Scan(value)
	if err != nil {
		err = new(apierror.ApiError).From(err)
		return
	}

	if v.Valid {
		*b, err = base64.StdEncoding.DecodeString(v.String)
		if err != nil {
			err = new(apierror.ApiError).From(err)
			return
		}
	}

	return
}

func (b Bytes) Value() (driver.Value, error) {
	return base64.StdEncoding.EncodeToString(b), nil
}

type Timestamp struct {
	time.Time
	Valid bool
}

func (t *Timestamp) SetTime(ti time.Time) {
	t.Time = ti
	t.Valid = true
}

func (t *Timestamp) Scan(value interface{}) (err error) {
	var v sql.NullInt64
	err = v.Scan(value)
	if err != nil {
		err = new(apierror.ApiError).From(err)
		return
	}

	if v.Valid {
		k := decimal.NewFromInt(1000)
		raw := decimal.NewFromInt(v.Int64)
		seconds := raw.Div(k).IntPart()
		nanoseconds := raw.Mod(k.Mul(k)).Mul(k).Mul(k).IntPart()
		t.SetTime(time.Unix(seconds, nanoseconds).In(Now().Location()))
	}

	return
}

func (t Timestamp) Value() (driver.Value, error) {
	if t.Valid {
		k := decimal.NewFromInt(1000)
		seconds := decimal.NewFromInt(t.Unix())
		nanoseconds := decimal.NewFromInt(int64(t.Nanosecond()))
		milliseconds := seconds.Mul(k).Add(
			nanoseconds.Div(k).Div(k),
		)

		return milliseconds.IntPart(), nil
	}

	return nil, nil
}
