package scoreboard

import (
	"context"
	"github.com/k-mistele/ccdc-scoreserver/lib/constants"
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

	//log.Debug("Stored serviceChecks")

}

// GET A LIST OF ServiceScoreCheck, IN THIS CASE ALL OF THEM
func GetAllServiceScoreChecks() (*[]ServiceScoreCheck, error ){

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

// GET A LIST OF THE MOST RECENT TIMES
func GetRecentScoreCheckTimes(numScoreChecks int, serviceName string) ([]string, error) {

	// SET UP A DATABASE CONNECTION
	client, ctx, err := database.GetClient()
	if err != nil {
		return []string{}, nil
	}
	defer client.Disconnect(*ctx)

	// GET THE COLLECTION
	collection := client.Database(database.Database).Collection(string(database.ServiceScoreCheck))

	// CREATE OPTIONS
	var scoreChecks []ServiceScoreCheck
	opts := options.Find()
	opts.SetSort(bson.M{"time": -1})
	opts.SetLimit(int64(numScoreChecks))

	// RUN THE QUERY AND DECODE IT
	cursor, err := collection.Find(context.TODO(), bson.M{"servicename": serviceName}, opts)
	if err != nil {
		return []string{}, err
	}

	// DESERIALIZE IT
	err = cursor.All(context.TODO(), &scoreChecks)
	if err != nil {
		return []string{}, err
	}

	// FIND ALL THE TIMES
	times := make([]string, 0)
	for _, scoreCheck := range scoreChecks {
		t := time.Unix(scoreCheck.Time, 0).In(constants.ServerTime).Format("15:04")
		times = append(times, t)
	}

	return times, nil

}

// GET A LIST OF THE n MOST RECENT ServiceScoreCheck OBJECTS
func GetRecentScoreChecks(numScoreChecks int, serviceName string) (*[]ServiceScoreCheck, error) {

	// SET UP A DATABASE CONNECTION
	client, ctx, err := database.GetClient()
	if err != nil {
		return &[]ServiceScoreCheck{}, err
	}
	defer client.Disconnect(*ctx)

	// GET THE COLLECTION
	collection := client.Database(database.Database).Collection(string(database.ServiceScoreCheck))

	// CREATE OPTIONS
	var scoreChecks []ServiceScoreCheck
	opts := options.Find()
	opts.SetSort(bson.M{"time": -1})
	opts.SetLimit(int64(numScoreChecks))

	// RUN THE QUERY
	cursor, err := collection.Find(context.TODO(), bson.M{"servicename": serviceName}, opts)
	if err != nil {
		return &scoreChecks, err
	}

	// DESERIALIZE IT
	err = cursor.All(context.TODO(), &scoreChecks)
	if err != nil {
		return &scoreChecks, err
	}

	return &scoreChecks, nil

}

// CHECK TO SEE IF A service.Service IS STILL UP BY CHECKING ITS MOST RECENT ServiceScoreCheck
func ServiceIsUp(serviceName string) (bool, error) {
	// SET UP DATABASE CONNECTION
	client, ctx, err := database.GetClient()
	if err != nil {
		log.Criticalf("Failed to get database connection: %s", err)
		return false, err
	}
	defer client.Disconnect(*ctx)

	// GET THE COLLECTION
	collection := client.Database(database.Database).Collection(string(database.ServiceScoreCheck))

	// CREATE OPTIONS
	var result = ServiceScoreCheck{}
	opts := options.FindOne().SetSort(bson.M{"time": -1})

	// RUN THE QUERY AND DECODE IT
	err = collection.FindOne(context.TODO(), bson.M{"servicename": serviceName}, opts).Decode(&result)
	if err != nil {
		log.Criticalf("Failed to return a scoreboard check: %s", err)
		return false, err
	}

	return result.IsUp, nil
}