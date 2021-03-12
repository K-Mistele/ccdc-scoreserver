package models

import "github.com/labstack/echo/v4"

type LoginModel struct {
	Messages 		MessagesModel
}

func NewLoginModel(c *echo.Context) (LoginModel, error) {

	model := LoginModel{
		Messages: NewMessagesModel(c),
	}

	return model, nil
}