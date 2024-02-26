package web

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cast"
)

func App() *pocketbase.PocketBase {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.Use(LoadAuthContextFromCookie(app))
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))
		e.Router.POST("/login/", LoginHandler(app))
		e.Router.GET("/logout/", LogoutHandler)
		e.Router.GET("/", RootHandler)
		return nil
	})

	return app
}

func RootHandler(c echo.Context) error {
	isGuestRaw := c.Get(ContextAuthIsGuestKey)
	isGuest := cast.ToBool(isGuestRaw)

	if isGuest {
		return c.HTML(http.StatusOK, "not logged in")
	}

	return c.HTML(http.StatusOK, "logged in")
}
