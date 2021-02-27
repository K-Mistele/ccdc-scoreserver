package scoreboard

import (
	"context"
	"github.com/k-mistele/ccdc-scoreserver/lib/database"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

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
	var interfaceSlice []interface{} // CAN ONLY ADD A SLICE OF INTERFACE TO COLLECTION WITH INSERTMANY

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