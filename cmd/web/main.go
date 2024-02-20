package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tokens"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/spf13/cast"
)

func loadAuthContextFromCookie(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenCookie, err := c.Request().Cookie("token")
			if err != nil || tokenCookie.Value == "" {
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

func main() {
	app := pocketbase.New()

	// serves static files from the provided public dir (if exists)
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))
		return nil
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.Use(loadAuthContextFromCookie(app))
		return nil
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.POST("/login/", func(c echo.Context) error {
			email := c.FormValue("email")
			password := c.FormValue("password")

			authRecord, err := app.Dao().FindAuthRecordByEmail("users", email)
			if err != nil {
				return err
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
		})

		return nil
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/", func(c echo.Context) error {
			info := apis.RequestInfo(c)
			admin := info.Admin       // nil if not authenticated as admin
			record := info.AuthRecord // nil if not authenticated as regular auth record

			isGuest := admin == nil && record == nil

			if isGuest {
				return c.HTML(http.StatusOK, "not logged in")

			}

			return c.HTML(http.StatusOK, "logged in")

		})

		return nil
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/logout/", func(c echo.Context) error {
			c.SetCookie(&http.Cookie{
				Name:     "token",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				Secure:   true,
				HttpOnly: true,
			})

			return c.HTML(http.StatusOK, "logged out")
		})

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
