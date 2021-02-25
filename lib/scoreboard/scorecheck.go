package scoreboard

import (
	"github.com/k-mistele/ccdc-scoreserver/lib/service"
	"sync"
	"time"
)

type ScoreCheck struct {
	Time 		time.Time
	Scores		map[string]bool			// MAPS SERVICE NAME TO SCORE
}


// RUN A SCORE CHECK, PRIVATE
func (sb *Scoreboard) runScoreCheck() ScoreCheck {

	var sc ScoreCheck
	var wg sync.WaitGroup
	var s service.Service

	log.Infof("Running score check")

	// CREATE A SCORE CHECK
	sc = ScoreCheck{
		Time:   time.Now(),
		Scores: map[string]bool{},
	}

	// CREATE A WAITGROUP
	wg = sync.WaitGroup{}

	// KICK OFF SERVICE CHECKS
	for _, s = range sb.Services {
		wg.Add(1)
		go s.DispatchServiceCheck(&(sc.Scores), &wg)
	}

	// WAIT FOR SERVICE CHECKS TO FINISH
	wg.Wait()

	// RETURN IT WHEN WE'RE DONE
	return sc

}