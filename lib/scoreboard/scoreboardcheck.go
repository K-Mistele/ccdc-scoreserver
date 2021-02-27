package scoreboard

import (
	"context"
	"github.com/k-mistele/ccdc-scoreserver/lib/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type ScoreboardCheck struct {
	Time 		int64					`bson:"time"`// UNIX TIME STAMP
	Scores		map[string]bool			`bson:"scores"`// MAPS SERVICE NAME TO SCORE
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

// GET THE LATEST ScoreboardCheck
func GetLatestScoreboardCheck() (ScoreboardCheck, error) {

	// SET UP DATABASE CONNECTION
	client, ctx, err := database.GetClient()
	if err != nil {
		log.Criticalf("Failed to get database connection: %s", err)
		return ScoreboardCheck{}, err
	}
	defer client.Disconnect(*ctx)

	// GET THE COLLECTION
	collection := client.Database(database.Database).Collection(string(database.ScoreboardCheck))

	// CREATE OPTIONS
	var result = ScoreboardCheck{}
	opts := options.FindOne().SetSort(bson.M{"time": 1})

	// RUN THE QUERY AND DECODE IT
	err = collection.FindOne(context.TODO(), bson.M{}, opts).Decode(&result)
	if err != nil {
		log.Criticalf("Failed to return a scoreboard check: %s", err)
		return ScoreboardCheck{}, err
	}

	// LOG AND RETURN IT
	log.Infof("Got latest scoreboard check: %v", result)
	return result, nil


}


