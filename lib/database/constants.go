package database

// TODO CHANGE THIS WHEN WE'RE BACK IN DOCKER
const URI = "mongodb://ccdc-scoreserver-database:27017/"
const Database = "ccdc-scoreserver"

type Collection string

const (
	ScoreboardCheck		Collection = "scoreboardcheck"
	ServiceScoreCheck	Collection = "servicescorecheck"
)
