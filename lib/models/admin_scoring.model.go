package view_models

import (
	"github.com/k-mistele/ccdc-scoreserver/lib/scoreboard"
	"github.com/labstack/echo/v4"
)

// THE AdminScoringModel IS THE TYPE THAT WILL BE PIPELINED TO THE TEMPLATE FOR /admin/scoring
type AdminScoringModel struct {
	Messages 					MessagesModel
	TimeStartedAt 				string
	TimeFinishesAt 				string
	ScoreboardIsRunning			bool
}

// CONSTRUCT AND RETURN AN AdminScoringModel
func NewAdminScoringModel(sb *scoreboard.Scoreboard, c *echo.Context) (AdminScoringModel, error ){

	// GET TIME STARTED AT
	timeStarted, timeFinished := sb.GetStartEndTimes()

	model := AdminScoringModel {
		Messages: 				NewMessagesModel(c),
		TimeStartedAt:			timeStarted,
		TimeFinishesAt:  		timeFinished,
		ScoreboardIsRunning: 	sb.Running,
	}

	return model,  nil
}