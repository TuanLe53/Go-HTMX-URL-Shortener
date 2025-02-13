package middlewares

import (
	"log"
	"net/http"

	"github.com/TuanLe53/Go-HTMX-URL-Shortener/db/models"
	"github.com/TuanLe53/Go-HTMX-URL-Shortener/handlers"
	"github.com/TuanLe53/Go-HTMX-URL-Shortener/pkg/auth"
	"github.com/TuanLe53/Go-HTMX-URL-Shortener/templates/components"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func AuthenticateJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var token *jwt.Token

		// Retrieve the access token from the cookies
		accessCookie, err := c.Cookie("access")
		if err != nil {
			// No access token, check for the refresh token
			if err == http.ErrNoCookie {
				refreshCookie, err := c.Cookie("refresh")
				if err != nil {
					log.Printf("No refresh token found: %v", err)
					return handlers.Render(c, components.AccessDenied())
				}

				// Validate the refresh token
				token, err = auth.ValidateToken(refreshCookie.Value)
				if err != nil {
					log.Printf("Invalid refresh token: %v", err)
					return handlers.Render(c, components.AccessDenied())
				}

				// Create a new access token using the claims from the refresh token
				newAccessClaims := auth.CreateJWTClaims(token.Claims.(jwt.MapClaims)["id"].(string), 15)
				newAccessToken, err := auth.GenerateToken(newAccessClaims)
				if err != nil {
					log.Printf("Failed to generate new access token: %v", err)
					return handlers.Render(c, components.ErrorMessage("Session expired. Please log in again."))
				}

				// Create a new access cookie with the new access token
				newAccessCookie, err := handlers.CreateCookie("access", newAccessToken, 15)
				if err != nil {
					log.Printf("Failed to create new access cookie: %v", err)
					return handlers.Render(c, components.ErrorMessage("Session expired. Please log in again."))
				}

				// Set the new access cookie
				c.SetCookie(newAccessCookie)
			} else {
				log.Printf("No access token found and not a cookie error: %v", err)
				return handlers.Render(c, components.AccessDenied())
			}
		} else {
			// Validate the access token
			token, err = auth.ValidateToken(accessCookie.Value)
			if err != nil {
				log.Printf("Invalid access token: %v", err)
				return handlers.Render(c, components.AccessDenied())
			}
		}

		// Retrieve the claims from the validated token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return handlers.Render(c, components.AccessDenied())
		}

		// Store the claims in the context for the next handler
		c.Set("user", claims)

		return next(c)

	}
}

func IsURLOwner(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(jwt.MapClaims)
		if !ok {
			return handlers.Render(c, components.ErrorMessage("Something went wrong. Please try again later."))
		}
		user_id := user["id"].(string)

		url_short_code := c.Param("short_code")

		url, err := models.GetURLDetail(url_short_code)
		if err != nil {
			log.Println("Error getting short URL:", err)
			return handlers.Render(c, components.ErrorMessage("Something went wrong. Please try again later."))
		}
		if url == nil {
			return handlers.Render(c, components.ErrorMessage("URL not found"))
		}

		if user_id != url.User_ID.String() {
			return handlers.Render(c, components.AccessDenied())
		}

		return next(c)
	}
}
