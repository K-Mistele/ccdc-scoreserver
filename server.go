package main

import (
	"github.com/k-mistele/ccdc-scoreserver/lib/auth"
	"github.com/k-mistele/ccdc-scoreserver/lib/scoreboard"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	logging "github.com/op/go-logging"
	"html/template"
	"io"

	"os"
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

// THE SCOREBOARD
var sb scoreboard.Scoreboard

func main() {


	/////////////////////////////////////////////////////////////////////////
	// LOGGING
	/////////////////////////////////////////////////////////////////////////

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
	sb = scoreboard.NewScoreboard()

	// INITIALIZE THE APP, SETTING UP A DEFAULT ROUTE AND STATIC DIRECTORY
	e := echo.New()

	/////////////////////////////////////////////////////////////////////////
	// MIDDLEWARES
	/////////////////////////////////////////////////////////////////////////
	// RECOVER FROM PANICS
	e.Use(middleware.Recover())

	// LOG HTTP REQUESTS
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: " (${status}) ${method} ${uri} ${remote_ip}\n",
	}))

	// JWT
	//e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
	//	Skipper:       	middleware.DefaultSkipper,
	//	SigningMethod: 	middleware.AlgorithmHS256,
	//	ContextKey:    	"user",
	//	TokenLookup:   	"header:" + echo.HeaderAuthorization,
	//	AuthScheme:    	"Bearer",
	//	Claims:        	jwt.MapClaims{},
	//	SigningKey:		utils.GenerateSecureToken(256),
	//}))
	e.Renderer = t

	/////////////////////////////////////////////////////////////////////////
	// INITIAL APP SETUP
	/////////////////////////////////////////////////////////////////////////

	// CREATE AN INITIAL ADMINISTRATIVE USER
	err := auth.CreateInitialAdminUser()
	if err != nil {
		log.Criticalf("Failed to create initial admin user: %s", err)
	}
	/////////////////////////////////////////////////////////////////////////
	// ROUTES: FRONTEND
	/////////////////////////////////////////////////////////////////////////

	// DEFINE THE SCOREBOARD ROUTE - PUBLIC
	e.GET("/", index)
	e.GET("/login", loginPage)
	e.POST("/login", login)

	// CREATE ROUTERS FOR TEAMS
	blueTeamRouter := e.Group("")
	blackTeamRouter := e.Group("")

	// BLUE TEAM VIEW ROUTES
	blueTeamRouter.GET("/services", services)
	//blueTeamRouter.GET("/injects", injects)
	//blueTeamRouter.GET("/reports", breachReports)

	// ADMIN GROUP VIEW ROUTES
	blackTeamRouter.GET("/admin/services/configure", adminConfigureServices)
	blackTeamRouter.GET("/admin/services/add", adminAddServices)
	blackTeamRouter.GET("/admin/scoring", adminScoring)

	/////////////////////////////////////////////////////////////////////////
	// ROUTES - BACKEND
	/////////////////////////////////////////////////////////////////////////

	blackTeamRouter.GET("/servicecheck/:name/params", getServiceCheckParams)

	// CHANGE THE PASSWORD FOR A SERVICE
	blueTeamRouter.POST("/service/:name/password", updateServicePassword)

	// DELETE A SERVICE
	blackTeamRouter.DELETE("/service/:name", deleteService)
	blackTeamRouter.PATCH("/service/:name", updateService)
	blackTeamRouter.PUT("/service/:name", createService)

	// START SCORING
	blackTeamRouter.POST("/scoring/start", startScoring)
	blackTeamRouter.POST("/scoring/restart", restartScoring)
	blackTeamRouter.POST("/scoring/stop", stopScoring)

	// STATIC DIRECTORY
	e.Static("/assets", "assets")

	// START THE APP
	e.Logger.Fatal(e.Start(":8080"))

}
