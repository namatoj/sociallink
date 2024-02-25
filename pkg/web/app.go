package web

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/template"
	"github.com/spf13/cast"
)

func App() *pocketbase.PocketBase {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		registry := template.NewRegistry()

		e.Router.Use(LoadAuthContextFromCookie(app))
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))
		e.Router.POST("/login/", LoginHandler(app))
		e.Router.GET("/logout/", LogoutHandler)
		e.Router.GET("/", RootHandler(registry))
		return nil
	})

	return app
}

func RootHandler(registry *template.Registry) echo.HandlerFunc {
	return func(c echo.Context) error {
		isGuestRaw := c.Get(ContextAuthIsGuestKey)
		isGuest := cast.ToBool(isGuestRaw)

		html, err := registry.LoadFiles(
			"pkg/web/templates/layout.html",
			"pkg/web/templates/root.html",
		).Render(map[string]any{
			"isGuest": isGuest,
			"user":    c.Get(apis.ContextAuthRecordKey),
		})

		if err != nil {
			return c.HTML(http.StatusInternalServerError, err.Error())
		}

		return c.HTML(http.StatusOK, html)
	}
}
