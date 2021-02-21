package main

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

func main() {

	// INITIALIZE THE APP, SETTING UP A DEFAULT ROUTE AND STATIC DIRECTORY
	e := echo.New()
	e.GET("/", func (c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to Kyle's CCDC Score Server")
	})
	e.Static("/assets", "assets")
	e.Logger.Fatal(e.Start(":8080"))

}