package database

// TODO CHANGE THIS WHEN WE'RE BACK IN DOCKER
const URI = "mongodb://localhost:27017/"
const Database = "ccdc-scoreserver"

type Collection string

const (
	ScoreboardCheck		Collection = "scoreboardcheck"
	ServiceScoreCheck	Collection = "servicescorecheck"
)
