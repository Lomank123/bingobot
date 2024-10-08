package internal_utils

import (
	"log"

	"fmt"
	"os"
	"path/filepath"
)

// Returns the absolute path of the given environment file (envFile) in the Go module's
// root directory. It searches for the 'go.mod' file from the current working directory upwards
// and appends the envFile to the directory containing 'go.mod'.
// It panics if it fails to find the 'go.mod' file.
// Reference: https://github.com/joho/godotenv/issues/126#issuecomment-1474645022
func Dir(envFile string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			break
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			panic(fmt.Errorf("go.mod not found"))
		}
		currentDir = parent
	}

	return filepath.Join(currentDir, envFile)
}

// Fetches variable from env file.
// If no key found it sets the `defaultVal` instead.
// Also prints the warning to the console in such case.
func GetEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	log.Printf("WARNING: Default value '%s' is used for key: '%s'", defaultVal, key)
	return defaultVal
}
