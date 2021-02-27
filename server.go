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
	"strconv"
	"time"

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
		templates: template.Must(template.ParseGlob("templates/*.html")),
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

	// DEFINE THE SCOREBOARD
	e.GET("/", func (c echo.Context) error {

		// GET SERVER TIME
		tz, _ := time.LoadLocation("America/Chicago")

		// GET MOST RECENT SCOREBOARD CHECK
		sbc, err := scoreboard.GetLatestScoreboardCheck()
		if err != nil {
			sbc = scoreboard.ScoreboardCheck{}
		}

		// COUNT UP SERVICES AND DOWN SERVICES
		numberUpServices, numberDownServices := 0, 0
		for key, _ := range sbc.Scores {
			if sbc.Scores[key] == true {
				numberUpServices += 1
			} else {
				numberDownServices += 1
			}
		}
		upProgressBarWidth, downProgressBarWidth := 0, 0
		if len(sb.Services) != 0 {
			upProgressBarWidth = int((float32(numberUpServices) / float32(len(sb.Services))) * 100)
			downProgressBarWidth = int((float32(numberDownServices) / float32(len(sb.Services))) * 100)
		}

		log.Debugf("%d:%d", numberUpServices, numberDownServices)

		// GET TIME STARTED AT
		var timeStarted, timeFinished string
		if !sb.Running {
			timeStarted = "00:00"
			timeFinished = "00:00"
		} else {
			timeStarted = time.Unix(sb.TimeStarted, 0).In(tz).Format("15:04")
			timeFinished = time.Unix(sb.TimeFinishes, 0).In(tz).Format("15:04")
		}



		// DEFINITION OF DATA
		data := struct {
			Overview					struct {
				TimeStartedAt 				string
				TimeFinishesAt				string
				NumberUpServices 			int
				NumberDownServices 			int
				NumberTotalServices 		int
				UpProgressBarWidth			int
				DownProgressBarWidth		int
			}

			ScoreboardCheck 			scoreboard.ScoreboardCheck

		}{
			// INSTANTIATION OF DATA
			Overview: struct {
				TimeStartedAt 				string
				TimeFinishesAt				string
				NumberUpServices 			int
				NumberDownServices 			int
				NumberTotalServices			int
				UpProgressBarWidth			int
				DownProgressBarWidth		int
			}{
				TimeStartedAt: 				timeStarted,
				TimeFinishesAt: 			timeFinished,
				NumberUpServices:  			numberUpServices,
				NumberDownServices:  		numberDownServices,
				NumberTotalServices:  		len(sb.Services),
				UpProgressBarWidth:   		upProgressBarWidth,
				DownProgressBarWidth:       downProgressBarWidth,
			},

			ScoreboardCheck:  			sbc,

		}
		log.Debug(data)

		return c.Render(http.StatusOK, "index.html", data)
	})

	e.GET("/scoring/start", func (c echo.Context) error {

		var intervalInt, hoursInt, minutesInt int
		var err error

		interval, hours, minutes := c.QueryParam("interval"), c.QueryParam("hours"), c.QueryParam("minutes")
		if intervalInt, err = strconv.Atoi(interval); err != nil {
			intervalInt	 = 60
		}
		if hoursInt, err = strconv.Atoi(hours); err != nil {
			return c.String(http.StatusBadRequest, "You must specify a number of hours to run for!")
		}
		if minutesInt, err = strconv.Atoi(minutes); err != nil {
			minutesInt = 0
		}



		if !sb.Running {
			sb.StartScoring(time.Duration(intervalInt), time.Duration(hoursInt), time.Duration(minutesInt))
			return c.String(http.StatusOK, "Scoring started!")
		} else {
			return c.String(http.StatusForbidden, "Scoring has already been started!")
		}
	})

	e.GET("/scoring/restart", func (c echo.Context) error {
		var intervalInt, hoursInt, minutesInt int
		var err error

		interval, hours, minutes := c.QueryParam("interval"), c.QueryParam("hours"), c.QueryParam("minutes")
		if intervalInt, err = strconv.Atoi(interval); err != nil {
			intervalInt	 = 60
		}
		if hoursInt, err = strconv.Atoi(hours); err != nil {
			return c.String(http.StatusBadRequest, "You must specify a number of hours to run for!")
		}
		if minutesInt, err = strconv.Atoi(minutes); err != nil {
			minutesInt = 0
		}


		sb.RestartScoring(time.Duration(intervalInt), time.Duration(hoursInt), time.Duration(minutesInt))
		return c.String(http.StatusOK, "Scoring restarted!")
	})

	e.GET("/scoring/stop", func (c echo.Context) error {
		if err := sb.StopScoring(); err != nil {
			return c.String(http.StatusForbidden, "Unable to stop the scoreboard - it is not running!")
		} else {
			return c.String(http.StatusOK, "Scoring stopped!")
		}
	})

	e.GET("/test-mongo", func (c echo.Context) error {
		databases := database.ListDatabases()
		return c.String(http.StatusOK, fmt.Sprintf("%v", databases))
	})

	e.GET("/test-collections", func (c echo.Context) error {
		//scoreboard.GetAllServiceChecks()
		//scoreboard.GetAllScoreboardChecks()

		return c.String(http.StatusOK, "check logs")
	})

	e.Static("/assets", "assets")
	e.Logger.Fatal(e.Start(":8080"))

}
