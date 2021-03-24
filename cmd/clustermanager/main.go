package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(
		middleware.CORS(),
		middleware.GzipWithConfig(middleware.GzipConfig{
			Level: 5,
		}),
	)

	e.Use(middleware.Secure())
}
