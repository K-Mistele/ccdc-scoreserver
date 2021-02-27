package scoreboard

import (
	"context"
	"github.com/k-mistele/ccdc-scoreserver/lib/service"
	"github.com/k-mistele/ccdc-scoreserver/lib/database"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

type ScoreboardCheck struct {
	Time 		int64					// UNIX TIME STAMP
	Scores		map[string]bool			// MAPS SERVICE NAME TO SCORE
}

type ServiceScoreCheck struct {
	Time 		int64					// UNIX TIME STAMP
	IsUp		bool 					// SERVICE IS UP?
	ServiceName string					// NAME OF THE SERVICE
}

// STORE ScoreboardCheck TO MONGO
// GOROUTINE SO STORE NOTHING
func storeScoreboardCheck(sbc *ScoreboardCheck) {

	// CREATE A CLIENT
	client, err := mongo.NewClient(options.Client().ApplyURI(database.URI))
	if err != nil {
		//log.Fatal(err)
		log.Critical("Failed to store service score checks!")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// CONNECT
	err = client.Connect(ctx)
	if err != nil {
		//log.Fatalf("Error occurred while connecting: %s",
		log.Critical("Failed to store service score checks")
	}
	defer client.Disconnect(ctx)

	// GET THE SERVICE SCORE CHECK COLLECTION
	collection := client.Database(database.Database).Collection(string(database.ScoreboardCheck))
	if err != nil {
		log.Error("Unable to get collection %s", string(database.ServiceScoreCheck))
	}

	_, err = collection.InsertOne(ctx, *sbc)
	if err != nil {
		log.Errorf("Failed to insert score checks: %s", err)
	}

	log.Debug("Stored scoreboardChecks")

}

// STORE A LIST OF ServiceScoreCheck
// GOROUTING SO RETURN NOTHING
func storeServiceScoreChecks(sscs *[]ServiceScoreCheck) {

	var collection *mongo.Collection
	var interfaceSlice []interface{}		// CAN ONLY ADD A SLICE OF INTERFACE TO COLLECTION WITH INSERTMANY

	// COPY SERVICE SCORE CHECKS FROM THEIR SLICE TO THE INTERFACE SLICE
	interfaceSlice = make([]interface{}, len(*sscs))
	for i, scorecheck := range *sscs {
		interfaceSlice[i] = scorecheck
	}


	// CREATE A CLIENT
	client, err := mongo.NewClient(options.Client().ApplyURI(database.URI))
	if err != nil {
		//log.Fatal(err)
		log.Critical("Failed to store service score checks!")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// CONNECT
	err = client.Connect(ctx)
	if err != nil {
		//log.Fatalf("Error occurred while connecting: %s",
		log.Critical("Failed to store service score checks")
	}
	defer client.Disconnect(ctx)

	// GET THE SERVICE SCORE CHECK COLLECTION
	collection = client.Database(database.Database).Collection(string(database.ServiceScoreCheck))
	if err != nil {
		log.Error("Unable to get collection %s", string(database.ServiceScoreCheck))
	}


	_, err = collection.InsertMany(ctx, interfaceSlice)
	if err != nil {
		log.Errorf("Failed to insert score checks: %s", err)
	}

	log.Debug("Stored serviceChecks")

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