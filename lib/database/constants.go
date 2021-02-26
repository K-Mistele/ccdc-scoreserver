package database

const databaseName = "ccdc-scoreserver"

type Collection string

const (
	ScoreboardCheck		Collection = "ScoreboardCheck"
	ServiceScoreCheck	Collection = "serviceScoreCheck"
)
