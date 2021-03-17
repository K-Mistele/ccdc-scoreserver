package auth

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/k-mistele/ccdc-scoreserver/lib/database"
	"github.com/k-mistele/ccdc-scoreserver/lib/utils"
	"github.com/op/go-logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)
var log = logging.MustGetLogger("main")

type User struct {
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
	Team     string `json:"team"`
	UUID     string `json:"uuid"`
	Hash     string `json:"hash"`
}

// CREATE A NEW User - ABSTRACT HASHING THE PASSWORD AND UUID GENERATION
func NewUser(username string, admin bool, team Team, password string) (User, error){

	// HASH THEIR PASSWORD

	hashBytes, err  := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	// BUILD A USER
	user := User {
		Username: username,
		Admin:    admin,
		Team:     string(team),
		UUID:     uuid.New().String(),
		Hash:     string(hashBytes),
	}

	return user, nil
}

func CreateInitialAdminUser() error {
	client, ctx, err := database.GetClient()
	if err != nil {
		log.Errorf("error creating initial admin user: %s", err)
		return err
	}
	defer client.Disconnect(*ctx)
	collection := client.Database(database.Database).Collection(string(database.User))

	// DELETE EXISTING ADMINS
	opts := options.Delete()
	_, err = collection.DeleteMany(context.TODO(), bson.M{"username": "admin"}, opts)
	if err != nil {
		log.Errorf("Error deleting admin users: %s", err)
	}

	// CREATE THE INITIAL ADMIN USER
	newAdminPass := utils.GenerateSecureToken(14)
	admin, err := NewUser("admin", true, Black, newAdminPass)
	if err != nil {
		panic(fmt.Sprintf("Error creating initial administrative user: %s", err))
	}
	err = admin.Store()
	if err != nil {
		panic(fmt.Sprintf("Error storing initial administrator: %s", err))
	}
	log.Infof("Initial admin (Black team) user created: %s:%s", admin.Username, newAdminPass)
	return nil
}
// STORE A User
func (user *User) Store() error {

	client, ctx, err := database.GetClient()
	if err != nil {
		return err
	}
	defer client.Disconnect(*ctx)

	collection := client.Database(database.Database).Collection(string(database.User))

	// INSERT
	_, err = collection.InsertOne(context.TODO(), *user)
	return err

}

// RETRIEVE A User BY User.UUID
func (User) GetByUUID(uuid string) (*User, error) {

	client, ctx, err := database.GetClient()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(*ctx)

	collection := client.Database(database.Database).Collection(string(database.User))

	user := User{}
	opts := options.FindOne()

	err = collection.FindOne(context.TODO(), bson.M{"uuid": uuid}, opts).Decode(&user)
	return &user, err
}

// RETRIEVE A User BY User.Username
func GetUserByUsername(username string) (*User, error) {

	client, ctx, err := database.GetClient()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(*ctx)

	collection := client.Database(database.Database).Collection(string(database.User))

	user := User{}
	opts := options.FindOne()

	err = collection.FindOne(context.TODO(), bson.M{"username": username}, opts).Decode(&user)
	return &user, err
}

// GET A TOKEN FOR A User
func (user *User) GetToken() (string, error) {
	return NewJSONWebToken(user.Username, Team(user.Team), user.Admin, user.UUID)
}

// LOGIN: CHECK A USERNAME AND PASSWORD, RETURNING A USER, TOKEN, AND BOOL
func Login(username string, password string) (user *User, token string, ok bool) {


	// SET UP A DATABASE CONNECTION
	user, err := GetUserByUsername(username)
	if err != nil {
		return nil, "", false
	}

	// CHECK THE PASSWORD OF THE USER
	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))
	if err != nil {
		log.Error(err)
		return nil, "", false
	}

	// IF NO ERROR WAS RETURNED, WE'RE GOOD TO ISSUE A TOKEN
	token, err = user.GetToken()
	if err != nil {
		log.Error(err)
		return nil, "", false
	}

	return user, token, true

}

// RETRIEVE ALL User OBJECTS FROM THE DATABASE
func GetAllUsers() (*[]User, error) {

	// SET UP DATABASE CONNECTION
	client, ctx, err := database.GetClient()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(*ctx)

	// FETCH ALL USERS
	var users []User
	opts := options.Find()
	collection := client.Database(database.Database).Collection(string(database.User))
	cursor, err := collection.Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.TODO(), &users)
	if err != nil {
		return nil, err
	}
	return &users, nil

}

// HANDLE CHANGING A User'S PASSWORD, INCLUDING ABSTRACTING AWAY DATABASE OPERATIONS
func (user *User) ChangePassword(newPassword string ) error {

	// HASH THE PASSWORD
	hashBytes, err  := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// GET A SESSION
	client, ctx, err := database.GetClient()
	if err != nil {
		return err
	}
	defer client.Disconnect(*ctx)

	// SET UP A COLLECTION
	collection := client.Database(database.Database).Collection(string(database.User))
	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"username", user.Username}}
	update := bson.D{{"$set", bson.D{{"hash", string(hashBytes)}}}}

	_, err = collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return err
	}
	return nil

}

// CHANGE THE PASSWORD OF A User GIVEN A USERNAME INSTEAD OF A User OBJECT - MAY REQUIRE ONE LESS LOOKUP
func ChangeUserPassword(username string, newPassword string ) error {

	// HASH THE PASSWORD
	hashBytes, err  := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// GET A SESSION
	client, ctx, err := database.GetClient()
	if err != nil {
		return err
	}
	defer client.Disconnect(*ctx)

	// SET UP A COLLECTION
	collection := client.Database(database.Database).Collection(string(database.User))
	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"username", username}}
	update := bson.D{{"$set", bson.D{{"hash", string(hashBytes)}}}}

	_, err = collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return err
	}
	return nil

}