package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
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

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:  "user-ui/build",
		HTML5: true,
	}))

	// Set cache control to a year on static directory. Filenames are hashed so
	// they can safely be aggressively cached.
	static := e.Group("/static", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Cache-Control", "max-age=31536000")
			return next(c)
		}
	})
	static.Use(middleware.Static("user-ui/build/static"))

	port := ":9500"

	log.Println("starting server on port serving index", port)

	e.Logger.Fatal(e.Start(port))
}
