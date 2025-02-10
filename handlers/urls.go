package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/TuanLe53/Go-HTMX-URL-Shortener/db/models"
	"github.com/TuanLe53/Go-HTMX-URL-Shortener/templates/components"
	"github.com/TuanLe53/Go-HTMX-URL-Shortener/templates/pages"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type URLHandler struct{}

func (h URLHandler) UserURLs(c echo.Context) error {
	return Render(c, pages.UserUrls())
}

func (h URLHandler) ShortenURL(c echo.Context) error {
	user, ok := c.Get("user").(jwt.MapClaims)
	if !ok {
		return Render(c, components.ErrorMessage("Something went wrong. Please try again later."))
	}

	url := c.FormValue("url")

	expire, err := strconv.Atoi(c.FormValue("expired_at"))
	if err != nil {
		log.Println("Error parsing expired_at:", err)
		return Render(c, components.ErrorMessage("Something went wrong. Please try again later."))
	}
	expired_at := time.Now().Add(time.Hour * time.Duration(expire))

	user_id := user["id"].(string)
	created_by, err := uuid.Parse(user_id)
	if err != nil {
		log.Println("Error parsing UUID:", err)
		return Render(c, components.ErrorMessage("Something went wrong. Please try again later."))
	}

	short_code, err := GenerateShortCode(12)
	if err != nil {
		log.Println("Error generating short code:", err)
		return Render(c, components.ErrorMessage("Something went wrong. Please try again later."))
	}

	_, err = models.CreateShortURL(url, short_code, created_by, expired_at)
	if err != nil {
		log.Println("Error creating short URL:", err)
		return Render(c, components.ErrorMessage("Something went wrong. Please try again later."))
	}

	c.Response().Header().Set("hx-redirect", fmt.Sprintf("/url/%s", short_code))
	return c.NoContent(http.StatusSeeOther)
}

func (h URLHandler) URLDetail(c echo.Context) error {
	url_id := c.Param("url_id")

	url, err := models.GetShortURL(url_id)
	if err != nil {
		log.Println("Error getting short URL:", err)
		return Render(c, components.ErrorMessage("Something went wrong. Please try again later."))
	}
	if url == nil {
		return Render(c, components.ErrorMessage("URL not found"))
	}

	return Render(c, pages.URLDetail(url))
}
