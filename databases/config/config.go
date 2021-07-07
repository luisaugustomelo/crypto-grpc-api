package config

import (
	"os"
)

// Config for database
type Config struct {
	Username     string
	Password     string
	DatabaseName string
	URL          string
	Collection   string
}

// GetConfig returns hardcoded config
func GetConfig() *Config {
	db_port := os.Getenv("KLEVER_MONGODB_PORT")
	username := os.Getenv("KLEVER_MONGODB_USERNAME")
	password := os.Getenv("KLEVER_MONGODB_PASSWORD")
	database := os.Getenv("KLEVER_MONGODB_DATABASE")
	collection := os.Getenv("KLEVER_MONGODB_COLLECTION_TEST")

	// if flag.Lookup("test.v") == nil {
	// 	collection = os.Getenv("KLEVER_MONGODB_COLLECTION")
	// } else {
	// 	collection = os.Getenv("KLEVER_MONGODB_COLLECTION_TEST")
	// }

	return &Config{
		Username:     username,
		Password:     password,
		DatabaseName: database,
		URL:          "mongodb://mongodb:" + db_port,
		Collection:   collection,
	}
}
