package scoreboard

import (
	"github.com/k-mistele/ccdc-scoreserver/lib/service"
	logging "github.com/op/go-logging"
	"time"
)

var log = logging.MustGetLogger("main")

// STATES FOR THE SCORING ROUTINE
type state string

const (
	start     state = "started"
	pause     state = "paused"
	resume    state = "resumed"
	terminate state = "terminated"
)
func NewScoreboard() Scoreboard {
	return Scoreboard{
		ScoringTerminated: false,
		signaller: make(chan state),
	}
}
// TYPE Scoreboard
// DEAL WITH AGGREGATING SERVICES, AND RUNNING AND AGGREGATING SCORE CHECKS
type Scoreboard struct {
	Services  			[]service.Service
	signaller 			chan state
	ScoringTerminated 	bool
}

// START THE SCORING ROUTINE
func (sb *Scoreboard) StartScoring(scoringInterval time.Duration) {

	var curState state
	seconds := scoringInterval * time.Second
	log.Debugf("Scoring started on %s interval", seconds)
	// CREATE A TICKER
	ticker := time.NewTicker(seconds)
	sb.ScoringTerminated = false
	go func() {

		// INFINITE LOOP
		for {
			// DETERMINE WHAT WE'RE DOING
			select {

			// IF WE HAVE A SIGNAL
			case newState := <-sb.signaller:
				switch newState {

				case pause:
					// IF WE GET A PAUSE COMMAND, KILL THE TICKER AND KEEP RUNNING THE LOOP
					log.Debug("Scoring Paused!")
					curState = newState
					ticker.Stop()

				case resume:

					// RESUME ONLY IF THE STATE IS PAUSE
					if curState == pause {
						log.Debug("Scoring resumed!")
						ticker = time.NewTicker(scoringInterval * time.Second)
						curState = resume
					} else if curState == resume || curState == start {
						// DO NOTHING
						continue
					}
				case terminate:

					// KILL THE SCORING AND RETURN
					log.Debug("Terminating Scoring")
					ticker.Stop()
					return
				}

			// IF WE HAVE A TICK
			case <-ticker.C:
				// GO RUN THE SCORE CHECK
				go sb.runScoreCheck()

			// IF WE HAVE NEITHER
			default:
				continue
			}

		}
	}()

}

// PAUSE SCORING
func (sb *Scoreboard) PauseScoring() {

	sb.signaller <- pause
}

// RESUME SCORING
func (sb *Scoreboard) ResumeScoring() {

	sb.signaller <- resume
}

// STOP SCORING
func (sb *Scoreboard) TerminateScoring() {

	sb.signaller <- terminate
	sb.ScoringTerminated = true

}
