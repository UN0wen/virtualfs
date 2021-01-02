package utils

import "os"

// GetVar gets an environment variable with name name, and returns its value if its set
// If not, the function returns the default value
func GetVar(name string, _default string) string {
	env := os.Getenv(name)
	if env == "" {
		return _default
	}
	return env
}

// DBUser for the production/development database
var DBUser = GetVar("DB_USER", "postgres")

// DBPort for the production/development database
var DBPort = GetVar("DB_PORT", "5432")

// DBPassword for the production/development database
var DBPassword = GetVar("DB_PASSWORD", "postgres")

// DBHost for the production/development database
var DBHost = GetVar("DB_HOST", "localhost")

// DBName for the production/development database
var DBName = GetVar("DB_NAME", "virtualfs")

// ServerPort is the port the server listens on
var ServerPort = GetVar("SERVER_PORT", "8080")
