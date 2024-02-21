package main

import (
	"log"
	"os"

	"github.com/namatoj/sociallink/pkg/web"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.Use(web.LoadAuthContextFromCookie(app))
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))
		e.Router.POST("/login/", web.LoginHandler(app))
		e.Router.GET("/logout/", web.LogoutHandler)
		e.Router.GET("/", web.RootHandler)
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
