package config

import (
	"fmt"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

// GetConnectionString returns the database connection string.
// If flagURL is provided, it returns it directly. Otherwise, it attempts
// to load variables from a .env file and construct the string.
func GetConnectionString(flagURL string) (string, error) {
	if flagURL != "" {
		return flagURL, nil
	}

	// Try loading .env. Ignore error if it doesn't exist,
	// as variables might be set directly in the environment.
	_ = godotenv.Load()

	host := os.Getenv("PGSQL_HOST")
	port := os.Getenv("PGSQL_PORT")
	user := os.Getenv("PGSQL_USER")
	dbname := os.Getenv("PGSQL_DB")
	password := os.Getenv("PGSQL_PASSWORD")
	sslmode := os.Getenv("PGSQL_SSLMODE")

	if host == "" || port == "" || user == "" || dbname == "" {
		return "", fmt.Errorf("connection parameters are missing. Please provide --db flag or set PGSQL_HOST, PGSQL_PORT, PGSQL_USER, and PGSQL_DB in .env")
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s", host, port, user, dbname)
	if password != "" {
		connStr += fmt.Sprintf(" password=%s", password)
	}
	if sslmode != "" {
		connStr += fmt.Sprintf(" sslmode=%s", sslmode)
	}

	return connStr, nil
}

var (
	// password=value, password='quoted value' and ?password=value in URIs.
	reKV = regexp.MustCompile(`(password\s*=\s*)('(?:[^'\\]|\\.)*'|[^\s&]+)`)
	// postgres://user:pass@... and postgresql://user:pass@...
	reURI = regexp.MustCompile(`(postgres(?:ql)?://[^:@/]*:)[^@/]*@`)
)

// RedactConnectionString replaces the password in a connection string with asterisks.
// Supports both key-value (password=...) and URI (postgres://user:pass@...) formats.
func RedactConnectionString(connStr string) string {
	connStr = reKV.ReplaceAllString(connStr, "${1}***")
	return reURI.ReplaceAllString(connStr, "${1}***@")
}
