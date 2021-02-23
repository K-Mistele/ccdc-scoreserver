package main

import (
	"github.com/k-mistele/ccdc-scoreserver/lib/service"
	"github.com/labstack/echo/v4"
	logging "github.com/op/go-logging"
	"html/template"
	"io"
	"os"

	"net/http"
)

// LOGGING STUFF
var log = logging.MustGetLogger("main")

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{level:.4s} ▶ %{shortfunc} ▶ %{id:03x}%{color:reset} %{message}`,
)

// FOR TEMPLATING
type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {

	// ENABLE LOG LEVELS
	// For demo purposes, create two backend for os.Stderr.
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)

	// For messages written to backend2 we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	backend2Formatter := logging.NewBackendFormatter(backend2, format)

	// Only errors and more severe messages should be sent to backend1
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")

	// Set the backends to be used.
	logging.SetBackend(backend1Leveled, backend2Formatter)
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
			Host:              "10.0.1.56",
			Port:              9000,
			Name:              "MinIO S3 Bucket",
			TransportProtocol: "tcp",
			Username:          "minioadmin",
			Password:          "minioadmin",
			ServiceCheck:      service.MinioBucketCheck,
			ServiceCheckData:  nil,
			Points:            10,
		}

		_, _ = s.DispatchServiceCheck()


		return c.String(http.StatusOK, "Welcome to Kyle's CCDC Score Server")
	})
	e.GET("/url", func(c echo.Context) error {
		return c.String(http.StatusOK, "Check completed successfully!")
	})

	e.Static("/assets", "assets")
	e.Logger.Fatal(e.Start(":8080"))

}
