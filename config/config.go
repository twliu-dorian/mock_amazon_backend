package config

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"strconv"
	"time"

	"mock_amazon_backend/log"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"gopkg.in/guregu/null.v4"
)

type DatabaseConfig struct {
	Username     string
	Password     string
	Host         string
	Port         int64
	DatabaseName string
	Version      null.Int
}

type Config struct {
	TimeZone          *time.Location
	HTTPListenAddress string
	HTTPListenPort    int64
	Database          DatabaseConfig
	AuthPrivateKey    *ecdsa.PrivateKey
}

var Global *Config

func InitConfig() (err error) {
	Global = new(Config)

	// timeZone
	timeZone := os.Getenv("TIME_ZONE")
	if err = validation.Validate(&timeZone, validation.Required); err != nil {
		err = fmt.Errorf(`"TIME_ZONE" %w`, err)
		return
	}
	if Global.TimeZone, err = time.LoadLocation(timeZone); err != nil {
		err = fmt.Errorf(`error on parsing "TIME_ZONE": %w`, err)
		return
	}
	log.Info(log.LabelStartup, fmt.Sprintf("Loaded environment variable %s=%s", "TIME_ZONE", timeZone))

	// http listen address
	Global.HTTPListenAddress = os.Getenv("HTTP_LISTEN_ADDR")
	if err = validation.Validate(&Global.HTTPListenAddress, validation.Required, is.IPv4); err != nil {
		err = fmt.Errorf(`"HTTP_LISTEN_ADDR" %w`, err)
		return
	}
	log.Info(log.LabelStartup, fmt.Sprintf("Loaded environment variable %s=%s", "HTTP_LISTEN_ADDR", Global.HTTPListenAddress))

	// http listen port
	httpListenPortString := os.Getenv("HTTP_LISTEN_PORT")
	if err = validation.Validate(&httpListenPortString, validation.Required, is.Int); err != nil {
		err = fmt.Errorf(`"HTTP_LISTEN_PORT" %w`, err)
		return
	}
	if Global.HTTPListenPort, err = strconv.ParseInt(httpListenPortString, 10, 64); err != nil {
		err = fmt.Errorf(`error on parsing "HTTP_LISTEN_PORT": %w`, err)
		return
	}
	log.Info(log.LabelStartup, fmt.Sprintf("Loaded environment variable %s=%d", "HTTP_LISTEN_PORT", Global.HTTPListenPort))

	Global.Database.Host = os.Getenv("DB_HOST")
	if err = validation.Validate(&Global.Database.Host, validation.Required, is.Host); err != nil {
		err = fmt.Errorf(`"DB_HOST" %w`, err)
		return
	}
	log.Info(log.LabelStartup, fmt.Sprintf("Loaded environment variable %s=%s", "DB_HOST", Global.Database.Host))

	//
	dbPortString := os.Getenv("DB_PORT")
	if err = validation.Validate(&dbPortString, validation.Required, is.Port); err != nil {
		err = fmt.Errorf(`"DB_PORT" %w`, err)
		return
	}
	if Global.Database.Port, err = strconv.ParseInt(dbPortString, 10, 64); err != nil {
		err = fmt.Errorf(`error on parsing "DB_PORT": %w`, err)
		return
	}
	log.Info(log.LabelStartup, fmt.Sprintf("Loaded environment variable %s=%d", "DB_PORT", Global.Database.Port))

	//
	Global.Database.Username = os.Getenv("DB_USERNAME")
	if err = validation.Validate(&Global.Database.Username, validation.Required); err != nil {
		err = fmt.Errorf(`"DB_USERNAME" %w`, err)
		return
	}
	log.Info(log.LabelStartup, fmt.Sprintf("Loaded environment variable %s=%s", "DB_USERNAME", Global.Database.Username))

	//
	Global.Database.Password = os.Getenv("DB_PASSWORD")
	if err = validation.Validate(&Global.Database.Password, validation.Required); err != nil {
		err = fmt.Errorf(`"DB_PASSWORD" %w`, err)
		return
	}
	log.Info(log.LabelStartup, fmt.Sprintf("Loaded environment variable %s=%s", "DB_PASSWORD", "********"))

	//
	Global.Database.DatabaseName = os.Getenv("DB_NAME")
	if err = validation.Validate(&Global.Database.DatabaseName, validation.Required); err != nil {
		err = fmt.Errorf(`"DB_NAME" %w`, err)
		return
	}
	log.Info(log.LabelStartup, fmt.Sprintf("Loaded environment variable %s=%s", "DB_NAME", Global.Database.DatabaseName))

	// database schema version
	dbVersion := os.Getenv("DB_VERSION")
	if len(dbVersion) > 0 {
		if err = Global.Database.Version.UnmarshalText([]byte(dbVersion)); err != nil {
			err = fmt.Errorf(`error on parsing "DB_VERSION": %w`, err)
			return
		}
		if Global.Database.Version.Int64 < 0 {
			err = fmt.Errorf(`error on parsing "DB_VERSION": must be uint`)
			return
		}
		log.Info(log.LabelStartup, fmt.Sprintf("Loaded environment variable %s=%d", "DB_VERSION", Global.Database.Version.Int64))
	} else {
		Global.Database.Version.Valid = false
		log.Info(log.LabelStartup, fmt.Sprintf("Set global config %s=null", "DB_VERSION"))
	}
	
	return
}
