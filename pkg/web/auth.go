package web

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tokens"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/spf13/cast"
)

const (
	ContextAuthIsGuestKey string = "isGuest"
)

func LoginHandler(app *pocketbase.PocketBase) echo.HandlerFunc {
	return func(c echo.Context) error {
		email := c.FormValue("email")
		password := c.FormValue("password")

		authRecord, err := app.Dao().FindAuthRecordByEmail("users", email)
		if err != nil {
			return c.HTML(http.StatusUnauthorized, "401 - Unauthorized")
		}

		if !authRecord.ValidatePassword(password) {
			return c.HTML(http.StatusUnauthorized, "401 - Unauthorized")
		}

		token, tokenErr := tokens.NewRecordAuthToken(app, authRecord)
		if tokenErr != nil {
			return c.HTML(http.StatusInternalServerError, "500 - Internal Server Error")
		}

		c.SetCookie(&http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		return c.HTML(http.StatusOK, "/login/")
	}
}

func LogoutHandler(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
	})

	return c.HTML(http.StatusOK, "logged out")
}

func LoadAuthContextFromCookie(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenCookie, err := c.Request().Cookie("token")
			if err != nil || tokenCookie.Value == "" {
				c.Set(ContextAuthIsGuestKey, true)
				return next(c) // no token cookie
			}

			token := tokenCookie.Value

			claims, _ := security.ParseUnverifiedJWT(token)
			tokenType := cast.ToString(claims["type"])

			switch tokenType {
			case tokens.TypeAdmin:
				admin, err := app.Dao().FindAdminByToken(
					token,
					app.Settings().AdminAuthToken.Secret,
				)
				if err == nil && admin != nil {
					// "authenticate" the admin
					c.Set(apis.ContextAdminKey, admin)
				}
			case tokens.TypeAuthRecord:
				record, err := app.Dao().FindAuthRecordByToken(
					token,
					app.Settings().RecordAuthToken.Secret,
				)
				if err == nil && record != nil {
					// "authenticate" the app user
					c.Set(apis.ContextAuthRecordKey, record)
				}
			}

			return next(c)
		}
	}
}
