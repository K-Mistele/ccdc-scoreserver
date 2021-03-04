package view_models

import (
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
func NewAdminServicesCreateModel(c *echo.Context) (AdminServicesCreateModel, error) {

	log.Debug("Creating /admin/services/create model")
	model := AdminServicesCreateModel{
		Messages: 					NewMessagesModel(c),
		AvailableServiceChecks: 	[]string{},
	}
	log.Debug(model.Messages)

	// BUILD THE LIST OF AVAILABLE SERVICE CHECKS
	for key := range service.ServiceChecks {
		model.AvailableServiceChecks = append(model.AvailableServiceChecks, key)
	}

	return model, nil
}