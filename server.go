package main

import (
	"fmt"
	"github.com/k-mistele/ccdc-scoreserver/lib/service"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"log"
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

		var s = service.Service{
			Host:              "127.0.0.1",
			Port:              8081,
			Name:              "Score Server HTTP",
			TransportProtocol: "tcp",
			Username:          "",
			Password:          "",
			ServiceCheck:      service.HTTPGetStatusCodeCheck,
			ServiceCheckData:  make(map[string]string),
			Points:            10,
		}
		s.ServiceCheckData["url"] = fmt.Sprintf("http://%s:%d/url", s.Host, s.Port)
		s.ServiceCheckData["expectedContent"] = "failed"


		isAlive, err := s.DispatchServiceCheck()
		if err != nil {
			log.Printf("Failed while trying to perform check on %s\n", s.Name)
		}

		log.Printf("Service %s is alive? %t", s.Name, isAlive)
		return c.String(http.StatusOK, "Welcome to Kyle's CCDC Score Server")
	})
	e.GET("/url", func(c echo.Context) error {
		return c.String(http.StatusOK, "Check completed successfully!")
	})

	e.Static("/assets", "assets")
	e.Logger.Fatal(e.Start(":8080"))

}
