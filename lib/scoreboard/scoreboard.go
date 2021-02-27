package scoreboard

import (
	"context"
	"errors"
	"github.com/k-mistele/ccdc-scoreserver/lib/database"
	"github.com/k-mistele/ccdc-scoreserver/lib/service"
	logging "github.com/op/go-logging"
	"time"
)

var log = logging.MustGetLogger("main")

// STATES FOR THE SCORING ROUTINE

func NewScoreboard() Scoreboard {
	return Scoreboard{
		stopSignaller: 			make(chan bool),
		Services: 					make([]service.Service, 0),
		TimeStarted: 				0,
		TimeFinishes: 				0,
		Running: 					false,
	}
}
// TYPE Scoreboard
// DEAL WITH AGGREGATING SERVICES, AND RUNNING AND AGGREGATING SCORE CHECKS
type Scoreboard struct {
	Services  			[]service.Service
	stopSignaller 	chan bool
	TimeStarted 		int64 		// UNIX TIMESTAMP
	TimeFinishes 		int64		// UNIX TIMESTAMP
	Running				bool		// HAS SCORING STARTED
}

func (sb* Scoreboard) ClearScores() {

	// SET UP A CLIENT
	client, ctx, err := database.GetClient()
	defer client.Disconnect(*ctx)

	db := client.Database(database.Database)
	serviceChecks := db.Collection(string(database.ServiceScoreCheck))
	scoreboardChecks := db.Collection(string(database.ScoreboardCheck))

	if err = serviceChecks.Drop(context.TODO()); err != nil {
		log.Error("Failed to drop service score checks collection: %s", err)
	}
	if err = scoreboardChecks.Drop(context.TODO()); err != nil {
		log.Error("Failed to drop scoreboard checks collection: %s", err)
	}
}

// START THE SCORING ROUTINE
// SCORING INTERVAL IS IN SECONDS
// SCORING DURATION IS hourDuration + minuteDuration
func (sb *Scoreboard) StartScoring(scoringInterval time.Duration, hourDuration time.Duration, minuteDuration time.Duration) error {

	if sb.Running {
		return errors.New("cannot start the scoreboard when it is already running")
	}
	seconds := scoringInterval * time.Second
	log.Infof("Scoring started on %s interval", seconds)

	// CREATE A TICKER
	ticker := time.NewTicker(seconds)

	curTime := time.Now().UTC()
	finishTime := curTime.Add(hourDuration * time.Hour).Add(minuteDuration * time.Minute)

	sb.TimeStarted = curTime.Unix()
	sb.TimeFinishes = finishTime.Unix()
	sb.Running = true

	go func() {

		// INFINITE LOOP
		for {
			// DETERMINE WHAT WE'RE DOING
			select {

			// IF WE HAVE A SIGNAL
			case <-sb.stopSignaller:
				ticker.Stop()
				sb.Running = false
				return

			// IF WE HAVE A TICK
			case <-ticker.C:

				// CHECK IF IT'S TIME TO STOP
				if time.Now().Unix() > sb.TimeFinishes {
					ticker.Stop()
					sb.Running = false
					log.Info("It's time to stop scoring! Stopping Scoring.")
					return
				}
				// GO RUN THE SCORE CHECK
				go sb.runScoreCheck()
			}
		}
	}()

	return nil

}

func (sb *Scoreboard) RestartScoring(scoringInterval time.Duration, hourDuration time.Duration, minuteDuration time.Duration) {

	// SEND THE RESTART SIGNALLER
	log.Info("Clearing database and restarting scoring")
	go func(){

		// IF RUNNING, SIGNAL TO THE ROUTINE WE'RE RESTARTING
		if sb.Running {
			sb.stopSignaller <- true
		}

		// CLEAR SCORES & START SCORING
		sb.ClearScores()
		sb.StartScoring(scoringInterval, hourDuration, minuteDuration)
	}()

}

func (sb *Scoreboard) StopScoring() error {

	log.Debug("Stopping Scoring")
	if sb.Running {
		go func() {
			sb.stopSignaller <- true
		}()
		return nil
	} else {
		return errors.New("cannot stop the scoreboard - it is not running")
	}
}
