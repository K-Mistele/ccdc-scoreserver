package scoreboard

import (
	"github.com/k-mistele/ccdc-scoreserver/lib/service"
	"sync"
	"time"
)

type ScoreboardCheck struct {
	Time 		int64					`bson:"time"`// UNIX TIME STAMP
	Scores		map[string]bool			`bson:"scores"`// MAPS SERVICE NAME TO SCORE
}

type ServiceScoreCheck struct {
	Time 		int64					`bson:"time"`// UNIX TIME STAMP
	IsUp		bool 					`bson:"isup"`// SERVICE IS UP?
	ServiceName string					`bson:"servicename"`// NAME OF THE SERVICE
}


// RUN A SCORE CHECK, PRIVATE
// GOROUTINE, SHOULD NOT RETURN ANYTHING
func (sb *Scoreboard) runScoreCheck() {

	var sbc ScoreboardCheck
	var wg sync.WaitGroup
	var s service.Service
	var curTime int64
	curTime = time.Now().UTC().Unix()
	var serviceScoreChecks []ServiceScoreCheck

	log.Infof("Running score check")

	// CREATE A ScoreboardCheck
	sbc = ScoreboardCheck{
		Time:   curTime,
		Scores: map[string]bool{},
	}
	log.Infof("Created Scoreboard check")

	// CREATE A WAITGROUP
	wg = sync.WaitGroup{}

	// KICK OFF SERVICE CHECKS

	for _, s = range sb.Services {
		wg.Add(1)
		go s.DispatchServiceCheck(&(sbc.Scores), &wg)
	}

	// WAIT FOR SERVICE CHECKS TO FINISH

	wg.Wait()

	// FOR EACH SERVICE IN THE SCOREBOARD CHECK, ADD A SERVICE
	serviceScoreChecks = make([]ServiceScoreCheck, len(sbc.Scores))
	i := 0
	for serviceName, isUp := range sbc.Scores {
		serviceScoreChecks[i] = ServiceScoreCheck{
			Time: curTime,
			IsUp: isUp,
			ServiceName: serviceName,
		}
		i++
	}

	// THROW BOTH INTO MONGO
	log.Debug("Storing information in database")
	go storeScoreboardCheck(&sbc)
	go storeServiceScoreChecks(&serviceScoreChecks)



}