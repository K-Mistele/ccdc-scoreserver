package models

import (
	"github.com/k-mistele/ccdc-scoreserver/lib/scoreboard"
	"github.com/k-mistele/ccdc-scoreserver/lib/service"
	"github.com/labstack/echo/v4"
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
// ALSO CONTAINS 1 LIST OF ALL OF THE service.Service	- USED FOR BUILDING PASSWORD CHANGE MODALS
type ServicesModel struct {

	// THEORETICALLY THERE SHOULD ONLY BE 2 COLUMNS
	Columns 			[][]ServiceCard
	Services 			[]service.Service
	Messages 			MessagesModel
}

// BUILD A ServicesModel AND RETURN IT FOR A ROUTE
func NewServiceModel(sb *scoreboard.Scoreboard, c *echo.Context) (ServicesModel, error) {

	servicesModel := ServicesModel{
		Columns: 	[][]ServiceCard{{}, {}},
		Services: 	[]service.Service{},
		Messages: 	NewMessagesModel(c),
	}
	for i, s := range sb.Services {

		// BUILD A SERVICE CARD
		serviceIsUp, _ := scoreboard.ServiceIsUp(s.Name)
		card := ServiceCard {
			Host: 				s.Host,
			Port: 				s.Port,
			Name: 				s.Name,
			Username: 			s.Username,
			TransportProtocol: 	s.TransportProtocol,
			Status: 			serviceIsUp,

		}

		// ADD THE service.Service TO THE LIST OF service.Service IN ServicesModel
		servicesModel.Services = append(servicesModel.Services, s)

		// ADD THE CARD TO THE APPROPRIATE COLUMN
		// EVEN SERVICES GO IN THE FIRST COLUMN, SINCE WE'RE 0-INDEXED
		if i % 2 == 0 {
			servicesModel.Columns[0] = append(servicesModel.Columns[0], card)
		} else {
			servicesModel.Columns[1] = append(servicesModel.Columns[1], card)
		}
	}

	return servicesModel, nil
}