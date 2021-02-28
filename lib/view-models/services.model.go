package view_models

import (
	"github.com/k-mistele/ccdc-scoreserver/lib/scoreboard"
	"github.com/k-mistele/ccdc-scoreserver/lib/service"
)

// THE ServicesModel STRUCT IS THE STRUCT THAT WILL BE PIPELINED TO THE TEMPLATE
// CONTAINS 2 LISTS OF scoreboard.ServiceScoreCheck - ONE FOR EACH TEMPLATE COLUMN
type ServicesModel struct {

	// THEORETICALLY THERE SHOULD ONLY BE COLUMNTS
	Columns 			[][]service.Service
}

// BUILD A ServicesModel AND RETURN IT FOR A ROUTE
func NewServiceModel(sb *scoreboard.Scoreboard) (ServicesModel, error) {

	servicesModel := ServicesModel{
		Columns: 	[][]service.Service{{}, {}},
	}
	for i, service := range sb.Services {

		// EVEN SERVICES GO IN THE FIRST COLUMN, SINCE WE'RE 0-INDEXED
		if i % 2 == 0 {
			servicesModel.Columns[0] = append(servicesModel.Columns[0], service)
		} else {
			servicesModel.Columns[1] = append(servicesModel.Columns[1], service)
		}
	}

	return servicesModel, nil
}