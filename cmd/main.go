package main

import (
	"log"

	"github.com/TuanLe53/Go-HTMX-URL-Shortener/db"
	"github.com/TuanLe53/Go-HTMX-URL-Shortener/db/models"
	"github.com/TuanLe53/Go-HTMX-URL-Shortener/handlers"
	"github.com/TuanLe53/Go-HTMX-URL-Shortener/middlewares"
	"github.com/TuanLe53/Go-HTMX-URL-Shortener/templates/pages"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	app := echo.New()

	app.Static("/static", "assets")

	app.Use(middleware.Logger())
	app.Use(middleware.Recover())

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.InitDB()
	db := db.DB()
	db.AutoMigrate(&models.User{}, &models.URL{}, &models.URLClick{})

	app.GET("/", func(c echo.Context) error {
		return handlers.Render(c, pages.Home())
	})

	authHandler := handlers.AuthHandler{}
	urlHandler := handlers.URLHandler{}

	app.GET("/login", authHandler.LoginPage)
	app.POST("/login", authHandler.LoginUser)
	app.GET("/register", authHandler.RegisterPage)
	app.POST("/register", authHandler.RegisterUser)

	app.GET("/short/:short_code", urlHandler.GoToURL)

	authGroup := app.Group("")
	authGroup.Use(middlewares.AuthenticateJWT)
	authGroup.POST("/shorten", urlHandler.ShortenURL)

	authGroup.GET("/my/urls", urlHandler.UserURLs)

	ownerGroup := authGroup.Group("/my")
	ownerGroup.Use(middlewares.IsURLOwner)
	ownerGroup.GET("/url/:short_code", urlHandler.URLDetail)
	ownerGroup.DELETE("/url/:short_code", urlHandler.DeleteURL)

	app.Logger.Fatal(app.Start(":5050"))
}
