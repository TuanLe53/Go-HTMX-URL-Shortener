package handlers

import (
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
