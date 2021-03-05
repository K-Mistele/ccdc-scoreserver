package main

import (
	"fmt"
	"github.com/k-mistele/ccdc-scoreserver/lib/messages"
	"github.com/k-mistele/ccdc-scoreserver/lib/scoreboard"
	"github.com/k-mistele/ccdc-scoreserver/lib/service"
	view_models "github.com/k-mistele/ccdc-scoreserver/lib/view-models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	// CREATE A SCOREBOARD
	var sb = scoreboard.NewScoreboard()

	// INITIALIZE THE APP, SETTING UP A DEFAULT ROUTE AND STATIC DIRECTORY
	e := echo.New()
	e.Use(middleware.Recover())
	e.Renderer = t

	// DEFINE THE SCOREBOARD ROUTE
	e.GET("/", func(c echo.Context) error {

		var model view_models.IndexModel
		model, err := view_models.NewIndexModel(&sb, &c)
		if err != nil {
			log.Criticalf("Error building index model: %s", err)
		}

		return c.Render(http.StatusOK, "index.html", model)
	})

	// DEFINE THE SERVICES ROUTES
	e.GET("/services", func(c echo.Context) error {
		var model view_models.ServicesModel
		model, err := view_models.NewServiceModel(&sb, &c)
		if err != nil {
			log.Critical("Error building services model: %s", err)
		}

		return c.Render(http.StatusOK, "services.html", model)
	})

	// DEFINE THE ADMIN SERVICES ROUTE
	e.GET("/admin/services/configure", func(c echo.Context) error {
		var model view_models.AdminServiceConfigModel
		model, err := view_models.NewAdminServicesConfigModel(&sb, &c)
		if err != nil {
			log.Critical("Error building Admin Services Model: %s", err)
		}

		return c.Render(http.StatusOK, "admin_services_configure.html", model)
	})

	// DEFINE THE ROUTE FOR CREATING SERVICES (FOR ADMINS ONLY)
	e.GET("/admin/services/add", func (c echo.Context) error {

		model, err  := view_models.NewAdminServicesCreateModel(&c)
		if err != nil {
			messages.Set(c, messages.Error, fmt.Sprint(err))
		}
		return c.Render(http.StatusOK, "admin_services_add.html", model)
	})

	// DEFINE THE ROUTE FOR MANAGING SCORING FOR ADMINS
	e.GET("/admin/scoring", func ( c echo.Context) error {
		model, err := view_models.NewAdminScoringModel(&sb, &c)
		if err != nil {
			messages.Set(c, messages.Error, fmt.Sprint(err))
		}
		return c.Render(http.StatusOK, "admin_scoring.html", model)
	})

	// GET THE PARAMETERS IF ANY FOR A SERVICE
	e.GET("/servicecheck/:name/params", func (c echo.Context) error {

		serviceName := c.Param("name")
		params, err := service.GetServiceParams(serviceName)
		if err != nil {
			params = []string{}
		}
		return c.JSON(http.StatusOK, params)
	})

	// CHANGE THE PASSWORD FOR A SERVICE
	e.POST("/service/:name/password", func(c echo.Context) error {

		var s *service.Service

		serviceName := c.Param("name")
		s, err := sb.GetService(serviceName)

		if err != nil {
			messages.Set(c, messages.Error, "Password change failed - service not found!")
			return c.Redirect(http.StatusFound, "/services")
		}

		password, confirmPassword := c.FormValue("password"), c.FormValue("confirmPassword")
		if password != confirmPassword {
			messages.Set(c, messages.Error, "Password change failed - passwords must match!")
			return c.Redirect(http.StatusFound, "/services")
		}

		s.ChangePassword(password)
		log.Infof("Changed password for service %s", serviceName)

		// FLASH A MESSAGE
		messages.Set(c, messages.Success, "Password successfully changed!")
		return c.Redirect(http.StatusFound, "/services")
	})

	// DELETE A SERVICE
	e.DELETE("/service/:name", func(c echo.Context) error {

		serviceName := c.Param("name")
		log.Infof("Attempting to delete service %s", serviceName)
		err := sb.DeleteService(serviceName)
		if err != nil {
			log.Errorf("Unable to delete service %s: %s", serviceName, err)
			messages.Set(c, messages.Error, "Unable to delete service!")
		} else {
			log.Infof("Successfully deleted service %s", serviceName)
			messages.Set(c, messages.Success, "Successfully deleted the service!")
		}

		return c.String(http.StatusAccepted, "")

	})

	// UPDATE A SERVICE
	e.PATCH("/service/:name", func(c echo.Context) error {

		serviceName := c.Param("name")
		log.Infof("Attempting to update service %s", serviceName)

		// GET STUFF FROM THE FORM
		host, portStr := c.FormValue("host"), c.FormValue("port")
		transportProto, username := c.FormValue("transportProtocol"), c.FormValue("username")
		password := c.FormValue("password")

		// CONVERT PORT TO AN INTEGER
		port, err := strconv.Atoi(portStr)
		if err != nil {
			messages.Set(c, messages.Error, "Port must be a number!")
			return c.String(http.StatusBadRequest, "")
		}

		// UPDATE THE SERVICE
		err = sb.UpdateService(serviceName, host, port, transportProto, username, password)
		if err != nil {
			log.Errorf("Unable to update services %s: %s", serviceName, err)
			messages.Set(c, messages.Error, "Unable to update service!")
		} else {
			log.Infof("Successfully updated service %s", serviceName)
			messages.Set(c, messages.Success, "Successfully updated the service!")
		}

		return c.String(http.StatusAccepted, "")

	})

	// CREATE A SERVICE
	e.PUT("/service/:name", func (c echo.Context) error {

		// ENSURE WE'RE NOT ADDING A DUPLICATE
		serviceName := c.Param("name")
		_, err := sb.GetService(serviceName)
		if err == nil {
			messages.Set(c, messages.Error, "A service with that name already exists!")
			return c.String(http.StatusForbidden, "A service with that name already exists!")
		}

		port, err := strconv.Atoi(c.FormValue("port"))
		if err != nil {
			messages.Set(c, messages.Error, "Port must be an integer!")
			return c.String(http.StatusBadRequest, "Port must be an integer!");
		}

		// CREATE A SERVICE CHECK
		s := service.Service {
			Host: 				c.FormValue("host"),
			Port: 				port,
			Name:				c.FormValue("name"),
			TransportProtocol: 	c.FormValue("proto"),
			Username: 			c.FormValue("username"),
			Password: 			c.FormValue("password"),
			ServiceCheck:		service.ServiceChecks[c.FormValue("checkType")],
			Points: 			10,
			ServiceCheckData: 	map[string]interface{}{},
		}

		// ADD PARAMS
		for _, param := range service.ServiceCheckParams[c.FormValue("checkType")] {
			s.ServiceCheckData[param] = c.FormValue(param)
		}
		log.Infof("Creating new service %v", s)

		sb.Services = append(sb.Services, s)

		messages.Set(c, messages.Success, "Service created!")
		return c.String(http.StatusCreated, "Service created!")

	})

	// START SCORING
	e.POST("/scoring/start", func(c echo.Context) error {

		var intervalInt, hoursInt, minutesInt int
		var err error

		interval, hours, minutes := c.FormValue("interval"), c.FormValue("hours"), c.FormValue("minutes")
		if intervalInt, err = strconv.Atoi(interval); err != nil {
			intervalInt = 60
		}
		if hoursInt, err = strconv.Atoi(hours); err != nil {
			return c.String(http.StatusBadRequest, "You must specify a number of hours to run for!")
		}
		if minutesInt, err = strconv.Atoi(minutes); err != nil {
			minutesInt = 0
		}

		if !sb.Running {
			err = sb.StartScoring(time.Duration(intervalInt), time.Duration(hoursInt), time.Duration(minutesInt))

			if err != nil {
				messages.Set(c, messages.Error, fmt.Sprint(err))
			} else {
				messages.Set(c, messages.Success, "Scoring started!")
			}

		} else {
			messages.Set(c, messages.Error, "Scoring is already running!")
		}
		return c.Redirect(http.StatusFound, "/admin/scoring")
	})

	e.POST("/scoring/restart", func(c echo.Context) error {
		var intervalInt, hoursInt, minutesInt int
		var err error

		interval, hours, minutes := c.FormValue("interval"), c.FormValue("hours"), c.FormValue("minutes")
		if intervalInt, err = strconv.Atoi(interval); err != nil {
			intervalInt = 60
		}
		if hoursInt, err = strconv.Atoi(hours); err != nil {
			return c.String(http.StatusBadRequest, "You must specify a number of hours to run for!")
		}
		if minutesInt, err = strconv.Atoi(minutes); err != nil {
			minutesInt = 0
		}

		sb.RestartScoring(time.Duration(intervalInt), time.Duration(hoursInt), time.Duration(minutesInt))
		messages.Set(c, messages.Success, "Scoring restarted!")
		return c.Redirect(http.StatusFound, "/admin/scoring")
	})

	e.POST("/scoring/stop", func(c echo.Context) error {
		if err := sb.StopScoring(); err != nil {
			messages.Set(c, messages.Error, fmt.Sprint(err))
		} else {
			messages.Set(c, messages.Success, "Scoring stopped!")
		}
		return c.Redirect(http.StatusFound, "/admin/scoring")
	})

	e.Static("/assets", "assets")
	e.Logger.Fatal(e.Start(":8080"))

}
