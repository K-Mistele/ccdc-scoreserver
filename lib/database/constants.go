package database

const URI = "mongodb://ccdc-scoreserver-database:27017/"
const Database = "ccdc-scoreserver"

type Collection string

const (
	ScoreboardCheck		Collection = "ScoreboardCheck"
	ServiceScoreCheck	Collection = "serviceScoreCheck"
)
