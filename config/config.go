package config

import (
	"flag"
	"os"
)

// Config holds application configuration
type Config struct {
	DebugMode bool
}

// Global configuration instance
var AppConfig Config

// Initialize configuration from command line flags and environment variables
func Init() {
	// Command line flags
	debugFlag := flag.Bool("debug", false, "Enable debug mode for detailed output")
	flag.Parse()

	// Check environment variable
	debugEnv := os.Getenv("BLACKJACK_DEBUG")
	
	// Set debug mode (command line flag takes precedence)
	AppConfig.DebugMode = *debugFlag || (debugEnv == "true" || debugEnv == "1")
}

// IsDebugMode returns whether debug mode is enabled
func IsDebugMode() bool {
	return AppConfig.DebugMode
}