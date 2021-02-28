package view_models

import (
	"github.com/k-mistele/ccdc-scoreserver/lib/scoreboard"
)

// THE ServiceCard STRUCT IS THE TYPE THAT'LL BE PIPELINED TO EACH CARD
type ServiceCard struct {
	Host					string
	Port					int
	Name					string
	Username 				string
	TransportProtocol 		string
	Status					bool
}

// THE ServicesModel STRUCT IS THE STRUCT THAT WILL BE PIPELINED TO THE TEMPLATE
// CONTAINS 2 LISTS OF scoreboard.ServiceScoreCheck - ONE FOR EACH TEMPLATE COLUMN
type ServicesModel struct {

	// THEORETICALLY THERE SHOULD ONLY BE COLUMNTS
	Columns 			[][]ServiceCard
}

// BUILD A ServicesModel AND RETURN IT FOR A ROUTE
func NewServiceModel(sb *scoreboard.Scoreboard) (ServicesModel, error) {

	servicesModel := ServicesModel{
		Columns: 	[][]ServiceCard{{}, {}},
	}
	for i, s := range sb.Services {

		serviceIsUp, _ := scoreboard.ServiceIsUp(s.Name)
		card := ServiceCard {
			Host: 				s.Host,
			Port: 				s.Port,
			Name: 				s.Name,
			Username: 			s.Username,
			TransportProtocol: 	s.TransportProtocol,
			Status: 			serviceIsUp,

		}

		// EVEN SERVICES GO IN THE FIRST COLUMN, SINCE WE'RE 0-INDEXED
		if i % 2 == 0 {
			servicesModel.Columns[0] = append(servicesModel.Columns[0], card)
		} else {
			servicesModel.Columns[1] = append(servicesModel.Columns[1], card)
		}
	}

	return servicesModel, nil
}