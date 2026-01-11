package config

import (
	"strconv"
	"time"
)

// TelnetConfig holds configuration for Telnet connections
type TelnetConfig struct {
	Host           string        // OLT IP address
	Port           int           // Telnet port (usually 23)
	Username       string        // Telnet username
	Password       string        // Telnet password
	EnablePassword string        // Enable mode password
	Timeout        time.Duration // Command timeout
	ConnectTimeout time.Duration // Connection timeout
	ReadTimeout    time.Duration // Read timeout
	WriteTimeout   time.Duration // Write timeout
	RetryCount     int           // Number of retry attempts
	RetryDelay     time.Duration // Delay between retries
	PoolSize       int           // Connection pool size
	MaxIdleTime    time.Duration // Max idle time before closing connection
	PromptUser     string        // User mode prompt (e.g., "ZXAN>")
	PromptEnable   string        // Enable mode prompt (e.g., "ZXAN#")
	PromptConfig   string        // Config mode prompt (e.g., "ZXAN(config)#")
}

// LoadTelnetConfig loads telnet configuration from environment variables
func LoadTelnetConfig() *TelnetConfig {
	port, _ := strconv.Atoi(getEnv("TELNET_PORT", "23"))
	timeout, _ := strconv.Atoi(getEnv("TELNET_TIMEOUT", "30"))
	connectTimeout, _ := strconv.Atoi(getEnv("TELNET_CONNECT_TIMEOUT", "10"))
	readTimeout, _ := strconv.Atoi(getEnv("TELNET_READ_TIMEOUT", "30"))
	writeTimeout, _ := strconv.Atoi(getEnv("TELNET_WRITE_TIMEOUT", "10"))
	retryCount, _ := strconv.Atoi(getEnv("TELNET_RETRY_COUNT", "3"))
	retryDelay, _ := strconv.Atoi(getEnv("TELNET_RETRY_DELAY", "2"))
	poolSize, _ := strconv.Atoi(getEnv("TELNET_POOL_SIZE", "1"))
	maxIdleTime, _ := strconv.Atoi(getEnv("TELNET_MAX_IDLE_TIME", "300"))

	return &TelnetConfig{
		Host:           getEnv("TELNET_HOST", "136.1.1.100"),
		Port:           port,
		Username:       getEnv("TELNET_USERNAME", "admin"),
		Password:       getEnv("TELNET_PASSWORD", ""),
		EnablePassword: getEnv("TELNET_ENABLE_PASSWORD", ""),
		Timeout:        time.Duration(timeout) * time.Second,
		ConnectTimeout: time.Duration(connectTimeout) * time.Second,
		ReadTimeout:    time.Duration(readTimeout) * time.Second,
		WriteTimeout:   time.Duration(writeTimeout) * time.Second,
		RetryCount:     retryCount,
		RetryDelay:     time.Duration(retryDelay) * time.Second,
		PoolSize:       poolSize,
		MaxIdleTime:    time.Duration(maxIdleTime) * time.Second,
		PromptUser:     getEnv("TELNET_PROMPT_USER", "ZXAN>"),
		PromptEnable:   getEnv("TELNET_PROMPT_ENABLE", "ZXAN#"),
		PromptConfig:   getEnv("TELNET_PROMPT_CONFIG", "ZXAN(config)#"),
	}
}

// Validate validates telnet configuration
func (c *TelnetConfig) Validate() error {
	if c.Host == "" {
		return ErrInvalidConfig("telnet host is required")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return ErrInvalidConfig("telnet port must be between 1 and 65535")
	}
	if c.Username == "" {
		return ErrInvalidConfig("telnet username is required")
	}
	if c.Password == "" {
		return ErrInvalidConfig("telnet password is required")
	}
	if c.Timeout <= 0 {
		return ErrInvalidConfig("telnet timeout must be positive")
	}
	if c.PoolSize < 1 {
		return ErrInvalidConfig("telnet pool size must be at least 1")
	}
	return nil
}

// ErrInvalidConfig creates a configuration error
func ErrInvalidConfig(message string) error {
	return &ConfigError{Message: message}
}

// ConfigError represents a configuration error
type ConfigError struct {
	Message string
}

func (e *ConfigError) Error() string {
	return "telnet config error: " + e.Message
}
