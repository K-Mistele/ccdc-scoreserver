package main

import (
	"fmt"
	"github.com/k-mistele/ccdc-scoreserver/lib/scoreboard"
	"github.com/k-mistele/ccdc-scoreserver/lib/service"
	"github.com/k-mistele/ccdc-scoreserver/lib/database"
	"github.com/labstack/echo/v4"
	logging "github.com/op/go-logging"
	"html/template"
	"io"
	"os"

	"net/http"
)

// LOGGING STUFF
var log = logging.MustGetLogger("main")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{level:.4s} ▶ %{id:03x}%{color:reset} %{message}`,
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

	// START THE SCORING

	// SET UP TEMPLATING STUFF
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/views/*.html")),
	}

	var sb = scoreboard.NewScoreboard()

	var s1 = service.Service{
		Host:              "10.0.1.51",
		Port:              445,
		Name:              "Charizard-DC-SMB",
		TransportProtocol: "tcp",
		Username:          "Administrator",
		Password:          "Password1!",
		ServiceCheck:      service.SMBListSharesCheck,
		ServiceCheckData:  map[string]interface{}{},
		Points:            10,
	}

	var s2 = service.Service {
		Host:				"10.0.1.52",
		Port: 				5900,
		Name: 				"Squirtle-VNC",
		TransportProtocol:  "tcp",
		Username: 			"",
		Password:			"password",
		ServiceCheck:		service.VNCConnectCheck,
		ServiceCheckData:   nil,
		Points:				10,
	}

	sb.Services = append(sb.Services, s1)
	sb.Services = append(sb.Services, s2)


	// INITIALIZE THE APP, SETTING UP A DEFAULT ROUTE AND STATIC DIRECTORY
	e := echo.New()
	e.Renderer = t

	// DEFINE A ROUTE
	e.GET("/", func(c echo.Context) error {

		return c.String(http.StatusOK, "Welcome to Kyle's CCDC Score Server")
	})

	// TODO: THESE SHOULD REALLY BE POST, BUT FOR TESTING THEY'RE GET
	e.GET("/scoring/start", func(c echo.Context) error {
		if !sb.ScoringTerminated {
			sb.StartScoring(5)
			return c.String(http.StatusOK, "Scoring started!")
		} else {
			return c.String(http.StatusForbidden, "Scoring has already been terminated")
		}

	})
	e.GET("/scoring/pause", func(c echo.Context) error {
		if !sb.ScoringTerminated {
			log.Debug("Pausing scoring")
			sb.PauseScoring()
			return c.String(http.StatusOK, "Scoring paused!")
		} else {
			return c.String(http.StatusForbidden, "Scoring has already been terminated")
		}
	})
	e.GET("/scoring/resume", func (c echo.Context) error {
		if !sb.ScoringTerminated {
			log.Debug("Resuming scoring")
			sb.ResumeScoring()
			return c.String(http.StatusOK, "Scoring resumed!")
		} else {
			return c.String(http.StatusForbidden, "Scoring has already been terminated")
		}

	})
	e.GET("/scoring/terminate", func (c echo.Context) error {
		if !sb.ScoringTerminated {
			log.Debug("Terminating scoring")
			sb.TerminateScoring()
			return c.String(http.StatusOK, "Scoring terminated")
		} else {
			return c.String(http.StatusForbidden, "Scoring has already been terminated")
		}

	})
	e.GET("/test-mongo", func (c echo.Context) error {
		databases := database.ListDatabases()
		return c.String(http.StatusOK, fmt.Sprintf("%v", databases))
	})

	e.Static("/assets", "assets")
	e.Logger.Fatal(e.Start(":8080"))

}
