package database

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"time"

	"mock_amazon_backend/apierror"
	"mock_amazon_backend/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	migrateMySQL "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

const (
	dbTimeout      string = "10s"
	dbWriteTimeout string = "30s"
	dbReadTimeout  string = "30s"
)

// InitDB initializes DB and pings it to make sure it is connectable
func InitDB() (err error) {
	if config.Global == nil {
		panic(new(apierror.ApiError).FromMessage("NO CONFIGS"))
	}

	username := config.Global.Database.Username
	password := config.Global.Database.Password
	host := config.Global.Database.Host
	port := config.Global.Database.Port
	dbName := config.Global.Database.DatabaseName

	address := fmt.Sprintf("%s:%d", host, port)

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=%s&readTimeout=%s&writeTimeout=%s&parseTime=true&loc=%s&multiStatements=true",
		username,
		password,
		address,
		dbName,
		dbTimeout,
		dbReadTimeout,
		dbWriteTimeout,
		url.QueryEscape(config.Global.TimeZone.String()), // Ensure time.Time parsing
	)

	db, err = sqlx.Open("mysql", dsn)
	if err != nil {
		err = new(apierror.ApiError).From(err)
		return
	}

	driver, err := migrateMySQL.WithInstance(db.DB, new(migrateMySQL.Config))
	if err != nil {
		err = new(apierror.ApiError).From(err)
		return
	}

	m, err := migrate.NewWithDatabaseInstance("file://sqls", dbName, driver)
	if err != nil {
		err = new(apierror.ApiError).From(err)
		return
	}

	if config.Global.Database.Version.Valid {
		err = m.Migrate(uint(config.Global.Database.Version.Int64))
	} else {
		err = m.Up()
	}
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		err = new(apierror.ApiError).From(err)
		return
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	_, err = db.Exec("SET time_zone = ?", config.Global.TimeZone.String())
	// _, err = db.Exec("SET time_zone = UTC")
	if err != nil {
		err = new(apierror.ApiError).From(err)
		return
	}

	return
}

func DB() (d *sqlx.DB, err error) {
	if db == nil {
		return nil, new(apierror.ApiError).FromMessage("DB is uninitialized")
	}
	return db, nil
}

func ClearTransition(tx *sqlx.Tx) {
	rollbackRet := tx.Rollback()
	if rollbackRet != sql.ErrTxDone && rollbackRet != nil {
		panic(rollbackRet.Error())
	}
}
