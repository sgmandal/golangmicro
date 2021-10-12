package main

import (
	"myswaggerautomation/api"
	_ "myswaggerautomation/docs"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
		Validator: func(username, password string, c echo.Context) (bool, error) {
			if username == "foo" && password == "bar" {
				return true, nil
			}
			return false, nil
		},
	}))
	e.POST("/foobar", api.FooBarHandler)
	e.Logger.Fatal(e.Start(":1323"))
}
