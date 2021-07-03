package config

import "os"

// Config for database
type Config struct {
	Username     string
	Password     string
	DatabaseName string
	URL          string
}

// GetConfig returns hardcoded config
func GetConfig() *Config {
	db_port := os.Getenv("KLEVER_MONGODB_PORT")
	username := os.Getenv("KLEVER_MONGODB_USERNAME")
	password := os.Getenv("KLEVER_MONGODB_PASSWORD")

	return &Config{
		Username: username,
		Password: password,
		URL:      "mongodb://mongodb:" + db_port,
	}
}
