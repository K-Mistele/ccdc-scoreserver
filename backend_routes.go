package main

import (
	"fmt"
	"github.com/k-mistele/ccdc-scoreserver/lib/messages"
	"github.com/k-mistele/ccdc-scoreserver/lib/service"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

// ROUTE FOR /SERVICECHECK/:NAME/PARAMS
func getServiceCheckParams(c echo.Context) error {

	serviceName := c.Param("name")
	params, err := service.GetServiceParams(serviceName)
	if err != nil {
		params = []string{}
	}
	return c.JSON(http.StatusOK, params)
}

// ROUTE FOR POST /SERVICE/:NAME/PASSWORD
func updateServicePassword(c echo.Context) error {

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
}

// ROUTE FOR DELETE /SERVICE/:NAME
func deleteService(c echo.Context) error {

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

}

// ROUTE FOR PATCH /SERVICE/:NAME
func updateService(c echo.Context) error {

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

}

// ROUTE FOR PUT /SERVICE/:NAME
func createService(c echo.Context) error {

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

}

// ROUTE FOR POST /SCORING/START
func startScoring(c echo.Context) error {

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

	if minutesInt + hoursInt <= 0 || intervalInt <= 0 {
		messages.Set(c, messages.Error, "Invalid scoring options! Total time and interval must be > 0")
		return c.Redirect(http.StatusFound, "/admin/scoring")
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
}

// ROUTE FOR POST /SCORING/RESTART
func restartScoring(c echo.Context) error {
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

	if minutesInt + hoursInt <= 0 || intervalInt <= 0 {
		messages.Set(c, messages.Error, "Invalid scoring options! Total time and interval must be > 0")
		return c.Redirect(http.StatusFound, "/admin/scoring")
	}

	sb.RestartScoring(time.Duration(intervalInt), time.Duration(hoursInt), time.Duration(minutesInt))
	messages.Set(c, messages.Success, "Scoring restarted!")
	return c.Redirect(http.StatusFound, "/admin/scoring")
}

// ROUTE FOR POST /SCORING/STOP
func stopScoring(c echo.Context) error {
	if err := sb.StopScoring(); err != nil {
		messages.Set(c, messages.Error, fmt.Sprint(err))
	} else {
		messages.Set(c, messages.Success, "Scoring stopped!")
	}
	return c.Redirect(http.StatusFound, "/admin/scoring")
}