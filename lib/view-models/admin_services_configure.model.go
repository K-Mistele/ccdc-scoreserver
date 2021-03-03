package view_models

import (
	"github.com/k-mistele/ccdc-scoreserver/lib/scoreboard"
	"github.com/k-mistele/ccdc-scoreserver/lib/service"
	"github.com/labstack/echo/v4"
)

// THE AdminServiceConfigCard STRUCT IS THE TYPE THAT'LL BE PIPELINED TO EACH CARD
type AdminServiceConfigCard struct {
	Host					string
	Port					int
	Name					string
	Username 				string
	Password 				string
	TransportProtocol 		string
	Status					bool
}

// THE AdminServiceConfigModel STRUCT IS THE STRUCT THAT WILL BE PIPELINED TO THE TEMPLATE
// CONTAINS 2 LISTS OF scoreboard.ServiceScoreCheck - ONE FOR EACH TEMPLATE COLUMN
// ALSO CONTAINS 1 LIST OF ALL OF THE service.Service	- USED FOR BUILDING PASSWORD CHANGE MODALS
type AdminServiceConfigModel struct {

	// THEORETICALLY THERE SHOULD ONLY BE 2 COLUMNS
	Columns 			[][]AdminServiceConfigCard
	Services 			[]service.Service
	Messages 			MessagesModel
}

// BUILD AN AdminServiceConfigModel AND RETURN IT FOR A ROUTE
func NewAdminServicesConfigModel(sb *scoreboard.Scoreboard, c *echo.Context) (AdminServiceConfigModel, error) {

	adminServicesModel := AdminServiceConfigModel{
		Columns: 	[][]AdminServiceConfigCard{{}, {}},
		Services: 	[]service.Service{},
		Messages: 	NewMessagesModel(c),
	}
	for i, s := range sb.Services {

		// BUILD A SERVICE CARD
		serviceIsUp, _ := scoreboard.ServiceIsUp(s.Name)
		card := AdminServiceConfigCard{
			Host: 				s.Host,
			Port: 				s.Port,
			Name: 				s.Name,
			Username: 			s.Username,
			Password: 			s.Password,
			TransportProtocol: 	s.TransportProtocol,
			Status: 			serviceIsUp,

		}

		// ADD THE service.Service TO THE LIST OF service.Service IN ServicesModel
		adminServicesModel.Services = append(adminServicesModel.Services, s)

		// ADD THE CARD TO THE APPROPRIATE COLUMN
		// EVEN SERVICES GO IN THE FIRST COLUMN, SINCE WE'RE 0-INDEXED
		if i % 2 == 0 {
			adminServicesModel.Columns[0] = append(adminServicesModel.Columns[0], card)
		} else {
			adminServicesModel.Columns[1] = append(adminServicesModel.Columns[1], card)
		}
	}

	return adminServicesModel, nil
}