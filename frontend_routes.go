package main

import (
	"fmt"
	"github.com/k-mistele/ccdc-scoreserver/lib/messages"
	"github.com/k-mistele/ccdc-scoreserver/lib/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

// ROUTE FOR INDEX /
func index(c echo.Context) error {

	var model models.IndexModel
	model, err := models.NewIndexModel(&sb, &c)
	if err != nil {
		log.Criticalf("Error building index model: %s", err)
	}

	return c.Render(http.StatusOK, "index.html", model)
}

// ROUTE FOR LOGIN
func loginPage(c echo.Context) error {

	model, err  := models.NewLoginModel(&c)
	if err != nil {
		log.Criticalf("Error building login model", err)
	}

	return c.Render(http.StatusOK, "login.html", model)
}

// ROUTE FOR /BLUETEAM/SERVICES
func services(c echo.Context) error {
	var model models.ServicesModel
	model, err := models.NewServiceModel(&sb, &c)
	if err != nil {
		log.Critical("Error building services model: %s", err)
	}

	return c.Render(http.StatusOK, "services.html", model)
}

// ROUTE FOR /BLACKTEAM/SERVICES/CONFIGURE
func adminConfigureServices(c echo.Context) error {
	var model models.AdminServiceConfigModel
	model, err := models.NewAdminServicesConfigModel(&sb, &c)
	if err != nil {
		log.Critical("Error building Admin Services Model: %s", err)
	}

	return c.Render(http.StatusOK, "admin_services_configure.html", model)
}

// ROUTE FOR /BLACKTEAM/SERVICES/ADD
func adminAddServices(c echo.Context) error {

	model, err  := models.NewAdminServicesCreateModel(&c)
	if err != nil {
		messages.Set(c, messages.Error, fmt.Sprint(err))
	}
	return c.Render(http.StatusOK, "admin_services_add.html", model)
}

// ROUTE FOR /BLACKTEAM/SCORING
func adminScoring( c echo.Context) error {
	model, err := models.NewAdminScoringModel(&sb, &c)
	if err != nil {
		messages.Set(c, messages.Error, fmt.Sprint(err))
	}
	return c.Render(http.StatusOK, "admin_scoring.html", model)
}

// ROUTE FOR GET /BLACKTEAM/USERS/ADD
func adminAddUsers (c echo.Context) error {
	model, err := models.NewAdminUsersAddModel(&c)
	log.Debug(model)
	if err != nil {
		messages.Set(c, messages.Error, fmt.Sprint(err))
		log.Error(err)
	}
	return c.Render(http.StatusOK, "admin_users_add.html", model)
}

// ROUTE FOR GET /BLACKTEAM/USERS/CONFIGURE
func adminConfigureUsers (c echo.Context) error {
	model, err := models.NewAdminUserConfigModel(&c)
	if err != nil {
		messages.Set(c, messages.Error, fmt.Sprint(err))
		log.Error(err)
	}
	return c.Render(http.StatusOK, "admin_users_configure.html", model)
}