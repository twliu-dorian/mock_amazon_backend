package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mock_amazon_backend/apierror"
	"mock_amazon_backend/config"
	"mock_amazon_backend/database"
	"mock_amazon_backend/log"
	"mock_amazon_backend/user"
	"mock_amazon_backend/util"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.HideBanner = true

	err := config.InitConfig()
	if err != nil {
		log.Fatal(log.LabelStartup, "Failed to start!", err)
	}

	err = database.InitDB()
	if err != nil {
		log.Fatal(log.LabelStartup, "Failed to start. Database ", err)
	}

	util.SetCurrentTimeFunc(func() time.Time {
		return time.Now().In(config.Global.TimeZone)
	})

	e.HTTPErrorHandler = apierror.HTTPErrorHandler

	e.Use(middleware.Gzip())
	e.Use(middleware.Recover())
	e.Use(log.RequestLogger())
	// e.Use(auth.CheckAuth())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "This is mock amazon backend server")
	})

	var userController = new(user.Controller)
	users := e.Group("/users")
	users.POST("/signUp", userController.Create)
	users.POST("/login", userController.Login)

	go func() {
		e.Logger.Info(log.LabelStartup, "Started successfully")
		if err := e.Start(fmt.Sprintf("%s:%d", config.Global.HTTPListenAddress, config.Global.HTTPListenPort)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(log.LabelShutdown, "System is shutting down, because of ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	osCall := <-quit

	e.Logger.Info(log.LabelShutdown, "System is shutting down, because of system call: ", osCall)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(log.LabelShutdown, err)
	}

	e.Logger.Info(log.LabelShutdown, "System successfully shut down.")

}
