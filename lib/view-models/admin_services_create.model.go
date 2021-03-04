package view_models

import (
	"github.com/k-mistele/ccdc-scoreserver/lib/scoreboard"
	"github.com/k-mistele/ccdc-scoreserver/lib/service"
	"github.com/labstack/echo/v4"
)

// THE AdminServicesCreateModel STRUCT IS THE STRUCT THAT WILL BE PIPELINED TO THE TEMPLATE
// IT CONTAINS DATA ABOUT THE SERVICE CONFIGURATION OPTIONS THAT ARE AVAILABLE
type AdminServicesCreateModel struct {
	Messages 					MessagesModel
	AvailableServiceChecks 		[]string
}

// BUILD AN AdminServiceCreateModel AND RETURN IT FOR A ROUTE
func NewAdminServicesCreateModel(sb *scoreboard.Scoreboard, c *echo.Context) (AdminServicesCreateModel, error) {

	model := AdminServicesCreateModel{
		Messages: 					NewMessagesModel(c),
		AvailableServiceChecks: 	[]string{},
	}

	// BUILD THE LIST OF AVAILABLE SERVICE CHECKS
	for key := range service.ServiceChecks {
		model.AvailableServiceChecks = append(model.AvailableServiceChecks, key)
	}

	return model, nil
}