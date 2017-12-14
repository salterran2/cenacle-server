package main

import (
	"fmt"
	"net/http"

	cfg "github.com/habasement/cenacle-server/config/env"

	"github.com/labstack/echo"
)

var config cfg.Config

func init() {
	config = cfg.NewViperConfig()

	if config.GetBool(`debug`) {
		fmt.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(config.GetString("server.address")))

}
