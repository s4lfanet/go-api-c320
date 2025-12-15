package utils

// GetConfigPath function determines the configuration path based on the environment string.
// Local, heroku, and production environments map to specific config directories.
func GetConfigPath(configPath string) string {
	if configPath == "development" { // Development environment
		return "./config/config-dev"
	} else if configPath == "heroku" { // Heroku environment
		return "./config/config-heroku"
	} else if configPath == "production" { // Production environment
		return "./config/config-prod"
	} else { // Default fallback
		return "./config/cfg"
	}
}
