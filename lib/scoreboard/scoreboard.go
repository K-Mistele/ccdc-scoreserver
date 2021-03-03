package scoreboard

import (
	"context"
	"errors"
	"fmt"
	"github.com/k-mistele/ccdc-scoreserver/lib/database"
	"github.com/k-mistele/ccdc-scoreserver/lib/service"
	logging "github.com/op/go-logging"
	"sync"
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

// RETURN A POINTER TO A service.Service WITH THE MATCHING NAME IN THE SCOREBOARD
func (sb *Scoreboard) GetService(serviceName string) (*service.Service, error) {

	// FIND A SERVICE WITH THE MATCHING NAME AND RETURN A POINTER
	for i, _ := range sb.Services {
		if sb.Services[i].Name == serviceName {
			return &(sb.Services[i]), nil
		}
	}

	// IF NONE FOUND, RETURN AN ERROR
	return nil, errors.New(fmt.Sprintf("unable to find a service on the scoreboard with a name of %s", serviceName))
}

// DELETE A service.Service IN THE SCOREBOARD WITH THE MATCHING NAME, OR RETURN AN ERROR
func (sb *Scoreboard) DeleteService(serviceName string) error {

	// FIND A SERVICE WITH THE MATCHING NAME AND IF IT EXISTS DELETE IT
	for i , _:= range sb.Services {
		if sb.Services[i].Name == serviceName {

			copy(sb.Services[i:], sb.Services[i+1:])
			sb.Services = sb.Services[:len(sb.Services)-1]
			return nil
		}
	}

	return errors.New(fmt.Sprintf("unable to delete services %s from scoreboard - unable to find it", serviceName))
}

// UPDATE A service.Service IN THE SCOREBOARD WITH THE MATCHING NAME, OR RETURN AN ERROR
func (sb *Scoreboard) UpdateService(serviceName string,
	host string,
	port int,
	transportProto string,
	username string,
	password string) error {

	// GET A POINTER TO THE SERVICE
	s, err := sb.GetService(serviceName)
	if err != nil {
		return err
	}
	s.Host = host
	s.Port = port
	s.TransportProtocol = transportProto
	s.Username = username
	s.Password = password
	return nil
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
	//log.Infof("Created Scoreboard check")

	// CREATE A WAITGROUP
	wg = sync.WaitGroup{}

	// KICK OFF SERVICE CHECKS

	for _, s = range (*sb).Services {
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
		//log.Debugf("Adding service Score check %s: %v", serviceName, serviceScoreChecks[i])
		i++
	}

	// THROW BOTH INTO MONGO
	//log.Debug("Storing information in database")
	go storeScoreboardCheck(&sbc)
	go storeServiceScoreChecks(&serviceScoreChecks)

}

