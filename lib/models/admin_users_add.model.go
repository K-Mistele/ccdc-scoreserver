package models

import "github.com/labstack/echo/v4"

type AdminUsersAddModel struct {
	Messages 			MessagesModel
}

func NewAdminUsersAddModel (c *echo.Context) (AdminUsersAddModel, error) {

	model := AdminUsersAddModel{
		Messages: 		NewMessagesModel(c),
	}

	return model, nil
}
