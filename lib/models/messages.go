package models

import (
	"github.com/k-mistele/ccdc-scoreserver/lib/messages"
	"github.com/labstack/echo/v4"
)

type MessagesModel struct {
	Success 		[]string
	Error 			[]string
}

// CONSTRUCT AND RETURN A NEW MessagesModel CONTAINING LISTS OF SUCCESS AND ERROR MESSAGES
func NewMessagesModel(c *echo.Context) MessagesModel {

	return MessagesModel {
		Success: 		messages.Get(*c, messages.Success),
		Error: 			messages.Get(*c, messages.Error),
	}
}