package view_models

import (
	"github.com/k-mistele/ccdc-scoreserver/lib/scoreboard"
	"github.com/labstack/gommon/log"
	"time"
)

type Overview struct {
	TimeStartedAt 				string
	TimeFinishesAt				string
	NumberUpServices 			int
	NumberDownServices 			int
	NumberTotalServices			int
	UpProgressBarWidth			int
	DownProgressBarWidth		int
	TotalUpServices				int
	TotalDownServices 			int
}
type PieChart struct {
	TotalUpServices				int
	TotalDownServices			int
}
type IndexModel struct {
	Overview			Overview
	PieChart 			PieChart
	ScoreboardCheck     scoreboard.ScoreboardCheck
}

func NewIndexModel(sb *scoreboard.Scoreboard) (IndexModel, error) {

	var indexModel IndexModel
	var pieChart PieChart
	var overview Overview

	////////////////////////////////////////////////////////
	// BUILD THE OVERVIEW - BARS AT THE TOP
	///////////////////////////////////////////////////////

	// GET SERVER TIME
	tz, _ := time.LoadLocation("America/Chicago")

	// GET MOST RECENT SCOREBOARD CHECK
	sbc, err := scoreboard.GetLatestScoreboardCheck()
	if err != nil {
		sbc = scoreboard.ScoreboardCheck{}
	}

	// COUNT UP SERVICES AND DOWN SERVICES
	numberUpServices, numberDownServices := 0, 0
	for key, _ := range sbc.Scores {
		if sbc.Scores[key] == true {
			numberUpServices += 1
		} else {
			numberDownServices += 1
		}
	}

	upProgressBarWidth, downProgressBarWidth := 0, 0
	if len(sb.Services) != 0 {
		upProgressBarWidth = int((float32(numberUpServices) / float32(len(sb.Services))) * 100)
		downProgressBarWidth = int((float32(numberDownServices) / float32(len(sb.Services))) * 100)
	}

	log.Debugf("%d:%d", numberUpServices, numberDownServices)

	// GET TIME STARTED AT
	var timeStarted, timeFinished string
	if !sb.Running {
		timeStarted = "00:00"
		timeFinished = "00:00"
	} else {
		timeStarted = time.Unix(sb.TimeStarted, 0).In(tz).Format("15:04")
		timeFinished = time.Unix(sb.TimeFinishes, 0).In(tz).Format("15:04")
	}

	overview = Overview{
		TimeStartedAt: 				timeStarted,
		TimeFinishesAt: 			timeFinished,
		NumberUpServices:  			numberUpServices,
		NumberDownServices:  		numberDownServices,
		NumberTotalServices:  		len(sb.Services),
		UpProgressBarWidth:   		upProgressBarWidth,
		DownProgressBarWidth:       downProgressBarWidth,
	}
	////////////////////////////////////////////////////////
	// BUILD THE PIECHART
	///////////////////////////////////////////////////////


	pieChart = PieChart{

	}

	////////////////////////////////////////////////////////
	// BUILD THE INDEX MODEL - BARS AT THE TOP
	///////////////////////////////////////////////////////

	indexModel = IndexModel{
		Overview: overview,
		PieChart: pieChart,
		ScoreboardCheck: sbc,
	}

	return indexModel, nil
}