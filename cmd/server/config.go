package main

import (
	"os"
)

// Config holds the application configuration.
type Config struct {
	OpenAIKey string
}

// LoadConfig loads the application configuration from environment variables.
func LoadConfig() *Config {
	return &Config{
		OpenAIKey: os.Getenv("OPENAI_API_KEY"),
	}
}
