package handlers

import (
	"log"
	"net/http"

	"github.com/TuanLe53/Go-HTMX-URL-Shortener/db/models"
	"github.com/TuanLe53/Go-HTMX-URL-Shortener/pkg/auth"
	"github.com/TuanLe53/Go-HTMX-URL-Shortener/templates/components"
	"github.com/TuanLe53/Go-HTMX-URL-Shortener/templates/pages"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct{}

func (h AuthHandler) LoginPage(c echo.Context) error {
	return Render(c, pages.Login())
}

func (h AuthHandler) RegisterPage(c echo.Context) error {
	return Render(c, pages.Register())
}

func (h AuthHandler) RegisterUser(c echo.Context) error {
	email := c.FormValue("email")
	firstName := c.FormValue("first_name")
	lastName := c.FormValue("last_name")
	password := c.FormValue("password")

	if !IsValidEmail(email) {
		return Render(c, components.ErrorMessage("Invalid email"))
	}

	isUserExists, err := models.FindUserWithEmail(email)
	if err != nil {
		return Render(c, components.ErrorMessage(err.Error()))
	}
	if isUserExists != nil {
		return Render(c, components.ErrorMessage("Email already taken."))
	}

	hashedPW, err := auth.HashPw(password)
	if err != nil {
		return Render(c, components.ErrorMessage(err.Error()))
	}

	models.CreateUser(email, firstName, lastName, hashedPW)

	c.Response().Header().Set("hx-redirect", "/login")
	return c.NoContent(http.StatusSeeOther)
}

func (h AuthHandler) LoginUser(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	if !IsValidEmail(email) {
		return Render(c, components.ErrorMessage("Invalid email"))
	}

	isUserExists, err := models.FindUserWithEmail(email)
	if err != nil {
		return Render(c, components.ErrorMessage("An error occurred, please try again later."))
	}

	if isUserExists == nil {
		return Render(c, components.ErrorMessage("User with this email does not exist."))
	}

	hashedPW := isUserExists.Password

	err = auth.CheckPw(hashedPW, password)
	if err != nil {
		return Render(c, components.ErrorMessage("Incorrect password."))
	}

	accessClaims := auth.CreateJWTClaims(isUserExists.ID.String(), 15)
	token, err := auth.GenerateToken(accessClaims)
	if err != nil {
		log.Println("Error generating token:", err)
		return Render(c, components.ErrorMessage("An error occurred, please try again later."))
	}

	accessCookie, err := CreateCookie("access", token, 15)
	if err != nil {
		log.Println("Error creating cookie", err)
		return Render(c, components.ErrorMessage("An error occurred, please try again later."))
	}

	refreshClaims := auth.CreateJWTClaims(isUserExists.ID.String(), 60*24) // Refresh token for 1 day
	refreshToken, err := auth.GenerateToken(refreshClaims)
	if err != nil {
		log.Println("Error generating refresh token:", err)
		return Render(c, components.ErrorMessage("An error occurred, please try again later."))
	}

	refreshCookie, err := CreateCookie("refresh", refreshToken, 60*24*7)
	if err != nil {
		log.Println("Error creating refresh cookie", err)
		return Render(c, components.ErrorMessage("An error occurred, please try again later."))
	}

	c.SetCookie(accessCookie)
	c.SetCookie(refreshCookie)

	c.Response().Header().Set("hx-redirect", "/")
	return c.NoContent(http.StatusSeeOther)
}
