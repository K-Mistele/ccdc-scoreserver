package scoreboard

import (
	"context"
	"github.com/k-mistele/ccdc-scoreserver/lib/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)
type ServiceScoreCheck struct {
	Time 		int64					`bson:"time"`// UNIX TIME STAMP
	IsUp		bool 					`bson:"isup"`// SERVICE IS UP?
	ServiceName string					`bson:"servicename"`// NAME OF THE SERVICE
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

// GET A LIST OF ServiceScoreCheck, IN THIS CASE ALL OF THEM
func GetServiceScoreChecks() (*[]ServiceScoreCheck, error ){

	var results []ServiceScoreCheck

	// SET UP DATABASE CONNECTION
	client, ctx, err := database.GetClient()
	if err != nil {
		results = make([]ServiceScoreCheck, 0)
		return &results, err
	}
	defer client.Disconnect(*ctx)

	// GET THE COLLECTION
	collection := client.Database(database.Database).Collection(string(database.ServiceScoreCheck))

	// CREATE OPTIONS
	opts := options.Find()

	// RUN THE FIND QUERY
	cursor, err := collection.Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		results = make([]ServiceScoreCheck, 0)
		return &results, err
	}

	// DESERIALIZE THEM
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		results = make([]ServiceScoreCheck, 0)
		return &results, err
	}

	return &results, nil

}

