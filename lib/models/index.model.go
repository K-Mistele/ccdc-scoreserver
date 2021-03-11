package models

import (
	"github.com/k-mistele/ccdc-scoreserver/lib/scoreboard"
	"github.com/labstack/echo/v4"
	"github.com/op/go-logging"
	cmap "github.com/orcaman/concurrent-map"
	"sync"
)

var log = logging.MustGetLogger("main")
var numChecksToDisplay = 10

// GOROUTINE TO GET THE MOST RECENT SCORES FOR A scoreboard.ServiceScoreCheck
// AND ADD THEM TO A MAP. SHOULD BE RUN WITH A WAITGROUP
func getRecentScoreCheckResults (serviceName string, results *cmap.ConcurrentMap, wg *sync.WaitGroup) {

	// GET THE SERVICE CHECKS
	serviceChecks, err := scoreboard.GetRecentScoreChecks(numChecksToDisplay, serviceName)
	if err != nil {
		log.Error("Failed to fetch score checks for %s: %s", serviceName, err)
		serviceChecks = &[]scoreboard.ServiceScoreCheck{}
	}

	// GET THE BOOLS FROM EACH AND ADD THEM TO A LIST
	checks := make([]bool, 0)
	for _, serviceCheck := range *serviceChecks {
		checks = append(checks, serviceCheck.IsUp)
	}

	// ADD THE SLICE TO THE MAP
	(*results).Set(serviceName, checks)

	// NOTIFY THE WAITGROUP WE'RE DONE
	wg.Done()

}

// THE Overview HANDLES THE CARDS AT THE TOP = TIME, PROGRESS BARS
type Overview struct {
	TimeStartedAt 				string
	TimeFinishesAt				string
	NumberUpServices 			int
	NumberDownServices 			int
	NumberTotalServices			int
	UpProgressBarWidth			int
	DownProgressBarWidth		int

}

// THE PieChart HANDLES THE PIECHART ON THE PAGE
type PieChart struct {
	TotalUpServices				int
	TotalDownServices			int
}

// THE ScoreboardChart IS CONTAINS THE DATA FOR BUILDING THE ACTUAL SCOREBOARD
type ScoreboardChart struct {
	Times 						[]string
	Services 					map[string][]bool
}

// THE IndexModel IS THE TYPE THAT'S PIPELINED TO THE TEMPLATE
type IndexModel struct {
	Overview			Overview
	PieChart 			PieChart
	ScoreboardChart		ScoreboardChart
	ScoreboardCheck     scoreboard.ScoreboardCheck
	Messages 			MessagesModel
}

// BUILD AN IndexModel AND RETURN IT FOR A ROUTE
func NewIndexModel(sb *scoreboard.Scoreboard, c *echo.Context) (IndexModel, error) {

	var indexModel IndexModel
	var pieChart PieChart
	var overview Overview
	var scoreboardChart ScoreboardChart
	var serviceScoreChecks *[]scoreboard.ServiceScoreCheck

	////////////////////////////////////////////////////////
	// BUILD THE Overview - BARS AT THE TOP
	///////////////////////////////////////////////////////

	// GET MOST RECENT SCOREBOARD CHECK
	sbc, err := scoreboard.GetLatestScoreboardCheck()
	if err != nil {
		sbc = scoreboard.ScoreboardCheck{}
	}

	// COUNT UP SERVICES AND DOWN SERVICES
	numberUpServices, numberDownServices := 0, 0
	for _, isUp := range sbc.Scores {
		//log.Debug(sbc.Scores)
		if isUp {
			numberUpServices += 1
		} else {
			numberDownServices += 1
		}
	}

	upProgressBarWidth, downProgressBarWidth := 0, 0
	if len(sb.Services) != 0 {
		//log.Debugf("up: %d, down: %d, total: %d", numberUpServices, numberDownServices, len(sb.Services))
		upProgressBarWidth = int((float32(numberUpServices) / float32(len(sb.Services))) * 100)
		downProgressBarWidth = int((float32(numberDownServices) / float32(len(sb.Services))) * 100)
	}

	//log.Debugf("%d:%d", numberUpServices, numberDownServices)

	// GET TIME STARTED AT
	timeStarted, timeFinished := sb.GetStartEndTimes()

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
	// BUILD THE PieChart
	///////////////////////////////////////////////////////
	totalUpServices, totalDownServices := 0, 0
	serviceScoreChecks, err = scoreboard.GetAllServiceScoreChecks()
	if err != nil {
		log.Criticalf("Failed to get service score checks: %s", err)
	}

	for _, scoreCheck := range *serviceScoreChecks {

		if scoreCheck.IsUp {
			totalUpServices += 1
		} else {
			totalDownServices += 1
		}
	}

	pieChart = PieChart{
		TotalDownServices: totalDownServices,
		TotalUpServices: totalUpServices,
	}

	////////////////////////////////////////////////////////
	// BUILD THE ScoreboardChart
	///////////////////////////////////////////////////////


	// GET TIMES OF RECENT SCORE CHECKS AS HH:MM
	var recentScoreCheckTimes []string
	if len(sb.Services) != 0 {
		recentScoreCheckTimes, err = scoreboard.GetRecentScoreCheckTimes(numChecksToDisplay, (*sb).Services[0].Name)
		if err != nil {
			log.Errorf("Couldn't get recent score check times: %s", recentScoreCheckTimes)
		}
	} else {
		log.Debug("Couldn't get scoreboard service times")
		recentScoreCheckTimes = []string{}
	}


	// BUILD THE MAP OF serviceNames TO LISTS OF bools - WHETHER THEY WERE UP OR NOT.
	// USE A WAITGROUP FOR THIS SINCE WE CAN DO IT CONCURRENTLY

	// TODO MAKE THIS A CMAP AND THEN COPY IT INTO A REGULAR MAP
	services := make(map[string][]bool)
	concurrentMap := cmap.New()
	wg := sync.WaitGroup{}

	// KICK OFF THE QUERIES
	for _, s := range sb.Services {
		wg.Add(1)
		go getRecentScoreCheckResults(s.Name, &concurrentMap, &wg)
	}

	// WAIT FOR ALL QUERIES TO FINISH
	wg.Wait()

	// COPY ALL THE VALUES FROM THE CONCURRENT MAP INTO THE NORMAL ONE
	for _, key := range concurrentMap.Keys() {
		v, ok := concurrentMap.Get(key)
		if !ok {
			log.Errorf("Couldn't get value %s from concurrent map", key)
		}
		services[key] = v.([]bool)
	}

	scoreboardChart = ScoreboardChart{
		Times: 		recentScoreCheckTimes,
		Services:	services,
	}

	//log.Debug(scoreboardChart)

	////////////////////////////////////////////////////////
	// BUILD THE IndexModel MODEL - BARS AT THE TOP
	///////////////////////////////////////////////////////

	indexModel = IndexModel{
		Overview: overview,
		PieChart: pieChart,
		ScoreboardChart: scoreboardChart,
		ScoreboardCheck: sbc,
		Messages: NewMessagesModel(c),
	}

	return indexModel, nil
}