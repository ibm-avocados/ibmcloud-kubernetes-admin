package main

import (
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	//e := echo.New()
	//e.Use(
	//	middleware.CORS(),
	//	middleware.GzipWithConfig(middleware.GzipConfig{
	//		Level: 5,
	//	}),
	//)
	//
	//e.Use(middleware.Secure())
	//
	//e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
	//	Root:  "user-ui/build",
	//	HTML5: true,
	//}))
	//
	//// Set cache control to a year on static directory. Filenames are hashed so
	//// they can safely be aggressively cached.
	//static := e.Group("/static", func(next echo.HandlerFunc) echo.HandlerFunc {
	//	return func(c echo.Context) error {
	//		c.Response().Header().Set("Cache-Control", "max-age=31536000")
	//		return next(c)
	//	}
	//})
	//static.Use(middleware.Static("user-ui/build/static"))
	//
	//auth := e.Group("/auth")
	//
	//auth.GET("", server.AuthHandler)
	//auth.GET("/callback", server.AuthDoneHandler) //url/auth/callback
	//
	//api := e.Group("/api/v1")
	//api.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	//	Format: "method=${method}, uri=${uri}, status=${status}, time=${latency_human}\n",
	//}))
	//
	//api.GET("/login", server.LoginHandler)
	//api.GET("/user/info", server.UserInfoHandler)
	//
	//port := ":9090"
	//
	//log.Println("starting server on port serving index", port)
	//
	//e.Logger.Fatal(e.Start(port))
}
