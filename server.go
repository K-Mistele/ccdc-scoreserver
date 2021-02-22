package main

import (
	"github.com/k-mistele/ccdc-scoreserver/lib"
	"github.com/labstack/echo/v4"
	"log"
	"html/template"
	"io"
	"net/http"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {

	// SET UP TEMPLATING STUFF
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/views/*.html")),
	}

	// INITIALIZE THE APP, SETTING UP A DEFAULT ROUTE AND STATIC DIRECTORY
	e := echo.New()
	e.Renderer = t

	// DEFINE A ROUTE
	e.GET("/", func(c echo.Context) error {

		var s = lib.Service {
			Host: "127.0.0.1",
			Port: 8080,
			Name: "Score Server",
			TransportProtocol: "tcp",
			Username: "admin",
			Password: "admin",
			ServiceCheckType: "tcp",
			ServiceCheckData: nil,
		}

		isAlive, err := s.DispatchServiceCheck()
		if err != nil {
			log.Printf("Failed while trying to perform check on %s\n", s.Name)
		}

		log.Printf("Service %s is alive? %t", s.Name, isAlive)
		return c.String(http.StatusOK, "Welcome to Kyle's CCDC Score Server")
	})

	e.Static("/assets", "assets")
	e.Logger.Fatal(e.Start(":8080"))

}
