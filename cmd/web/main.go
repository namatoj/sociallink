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
)

func main() {
	app := pocketbase.New()

	// serves static files from the provided public dir (if exists)
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))
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

			return c.HTML(http.StatusOK, "/login")
		})

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
