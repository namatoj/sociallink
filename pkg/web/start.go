package web

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
)

func RootHandler(c echo.Context) error {
	info := apis.RequestInfo(c)
	admin := info.Admin       // nil if not authenticated as admin
	record := info.AuthRecord // nil if not authenticated as regular auth record

	isGuest := admin == nil && record == nil

	if isGuest {
		return c.HTML(http.StatusOK, "not logged in")

	}

	return c.HTML(http.StatusOK, "logged in")
}
