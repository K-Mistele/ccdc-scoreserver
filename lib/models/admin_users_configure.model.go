package models

import (
	"github.com/k-mistele/ccdc-scoreserver/lib/auth"
	"github.com/labstack/echo/v4"
)

// THE AdminConfigureUserCard STORES INFORMATION ABOUT USERS TO BE TEMPLATED INTO CARDS
type AdminConfigureUserCard struct{
	Username 			string
	Team 				string
	IsAdmin 			bool
	UUID 				string
}

// THE AdminUserConfigModel IS THE STRUCT THAT IS PIPELINED TO THE TEMPLATE FOR RENDERING
type AdminUserConfigModel struct {
	Messages 				MessagesModel
	Columns					[][]AdminConfigureUserCard
	Users 					[]AdminConfigureUserCard
}

// THE NewAdminUserConfigModel() FUNCTION WILL BUILD AND RETURN AN AdminUserConfigModel FOR PIPELINING TO THE TEMPLATE
func NewAdminUserConfigModel(c *echo.Context) (AdminUserConfigModel, error) {

	// CREATE THE MODEL
	model := AdminUserConfigModel{
		Messages: 			NewMessagesModel(c),
		Columns: 			[][]AdminConfigureUserCard{{}, {}},
		Users: 				[]AdminConfigureUserCard{},
	}

	// GET ALL USERS
	users, err := auth.GetAllUsers()
	if err != nil {
		log.Criticalf("Failed to get all users: %s", err)
		return AdminUserConfigModel{}, err
	}

	// PROCESS THE USERS
	for i, user := range *users {

		// BUILD A CARD
		card := AdminConfigureUserCard{
			Username: 		user.Username,
			Team:     		user.Team,
			IsAdmin:  		user.Admin,
			UUID:     		user.UUID,
		}

		// ADD IT TO THE LIST OF ALL USERS
		model.Users = append(model.Users, card)

		// ADD IT TO THE CORRECT COLUMN
		if i % 2 == 0 {
			model.Columns[0] = append(model.Columns[0], card)
		} else {
			model.Columns[1] = append(model.Columns[1], card)
		}
	}

	return model, nil
}