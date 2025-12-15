package redis

import (
	"os"
	"testing"

	"github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/config"
)

func TestNewRedisClient_FromConfig(t *testing.T) {
	// Clear environment variables to ensure we use config
	os.Unsetenv("APP_ENV")
	os.Unsetenv("REDIS_HOST")
	os.Unsetenv("REDIS_PORT")
	os.Unsetenv("REDIS_PASSWORD")
	os.Unsetenv("REDIS_DB")
	os.Unsetenv("REDIS_MIN_IDLE_CONNECTIONS")
	os.Unsetenv("REDIS_POOL_SIZE")
	os.Unsetenv("REDIS_POOL_TIMEOUT")

	cfg := &config.Config{
		RedisCfg: config.RedisConfig{
			Host:               "localhost",
			Port:               "6379",
			Password:           "testpass",
			DB:                 1,
			MinIdleConnections: 10,
			PoolSize:           100,
			PoolTimeout:        30,
		},
	}

	client := NewRedisClient(cfg)

	if client == nil {
		t.Error("Expected non-nil Redis client")
	}

	opts := client.Options()

	expectedAddr := "localhost:6379"
	if opts.Addr != expectedAddr {
		t.Errorf("Expected address %s, got %s", expectedAddr, opts.Addr)
	}

	if opts.Password != "testpass" {
		t.Errorf("Expected password 'testpass', got %s", opts.Password)
	}

	if opts.DB != 1 {
		t.Errorf("Expected DB 1, got %d", opts.DB)
	}

	if opts.MinIdleConns != 10 {
		t.Errorf("Expected MinIdleConns 10, got %d", opts.MinIdleConns)
	}

	if opts.PoolSize != 100 {
		t.Errorf("Expected PoolSize 100, got %d", opts.PoolSize)
	}
}

func TestNewRedisClient_FromEnvironment_Development(t *testing.T) {
	// Set environment variables
	os.Setenv("APP_ENV", "development")
	os.Setenv("REDIS_HOST", "redis-dev")
	os.Setenv("REDIS_PORT", "6380")
	os.Setenv("REDIS_PASSWORD", "devpass")
	os.Setenv("REDIS_DB", "2")
	os.Setenv("REDIS_MIN_IDLE_CONNECTIONS", "20")
	os.Setenv("REDIS_POOL_SIZE", "200")
	os.Setenv("REDIS_POOL_TIMEOUT", "60")

	defer func() {
		os.Unsetenv("APP_ENV")
		os.Unsetenv("REDIS_HOST")
		os.Unsetenv("REDIS_PORT")
		os.Unsetenv("REDIS_PASSWORD")
		os.Unsetenv("REDIS_DB")
		os.Unsetenv("REDIS_MIN_IDLE_CONNECTIONS")
		os.Unsetenv("REDIS_POOL_SIZE")
		os.Unsetenv("REDIS_POOL_TIMEOUT")
	}()

	cfg := &config.Config{
		RedisCfg: config.RedisConfig{
			Host: "should-be-ignored",
			Port: "9999",
		},
	}

	client := NewRedisClient(cfg)

	if client == nil {
		t.Error("Expected non-nil Redis client")
	}

	opts := client.Options()

	expectedAddr := "redis-dev:6380"
	if opts.Addr != expectedAddr {
		t.Errorf("Expected address %s, got %s", expectedAddr, opts.Addr)
	}

	if opts.Password != "devpass" {
		t.Errorf("Expected password 'devpass', got %s", opts.Password)
	}

	if opts.DB != 2 {
		t.Errorf("Expected DB 2, got %d", opts.DB)
	}

	if opts.MinIdleConns != 20 {
		t.Errorf("Expected MinIdleConns 20, got %d", opts.MinIdleConns)
	}

	if opts.PoolSize != 200 {
		t.Errorf("Expected PoolSize 200, got %d", opts.PoolSize)
	}
}

func TestNewRedisClient_FromEnvironment_Production(t *testing.T) {
	// Set environment variables for production
	os.Setenv("APP_ENV", "production")
	os.Setenv("REDIS_HOST", "redis-prod")
	os.Setenv("REDIS_PORT", "6381")
	os.Setenv("REDIS_PASSWORD", "prodpass")
	os.Setenv("REDIS_DB", "3")
	os.Setenv("REDIS_MIN_IDLE_CONNECTIONS", "30")
	os.Setenv("REDIS_POOL_SIZE", "300")
	os.Setenv("REDIS_POOL_TIMEOUT", "90")

	defer func() {
		os.Unsetenv("APP_ENV")
		os.Unsetenv("REDIS_HOST")
		os.Unsetenv("REDIS_PORT")
		os.Unsetenv("REDIS_PASSWORD")
		os.Unsetenv("REDIS_DB")
		os.Unsetenv("REDIS_MIN_IDLE_CONNECTIONS")
		os.Unsetenv("REDIS_POOL_SIZE")
		os.Unsetenv("REDIS_POOL_TIMEOUT")
	}()

	cfg := &config.Config{
		RedisCfg: config.RedisConfig{
			Host: "should-be-ignored",
			Port: "9999",
		},
	}

	client := NewRedisClient(cfg)

	if client == nil {
		t.Error("Expected non-nil Redis client")
	}

	opts := client.Options()

	expectedAddr := "redis-prod:6381"
	if opts.Addr != expectedAddr {
		t.Errorf("Expected address %s, got %s", expectedAddr, opts.Addr)
	}

	if opts.Password != "prodpass" {
		t.Errorf("Expected password 'prodpass', got %s", opts.Password)
	}

	if opts.DB != 3 {
		t.Errorf("Expected DB 3, got %d", opts.DB)
	}
}

func TestNewRedisClient_WithEmptyPassword(t *testing.T) {
	os.Unsetenv("APP_ENV")

	cfg := &config.Config{
		RedisCfg: config.RedisConfig{
			Host:     "localhost",
			Port:     "6379",
			Password: "",
			DB:       0,
		},
	}

	client := NewRedisClient(cfg)

	if client == nil {
		t.Error("Expected non-nil Redis client")
	}

	opts := client.Options()

	if opts.Password != "" {
		t.Errorf("Expected empty password, got %s", opts.Password)
	}
}

func TestNewRedisClient_InvalidEnvironmentIntegers(t *testing.T) {
	// Set environment with invalid integer values
	os.Setenv("APP_ENV", "development")
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "6379")
	os.Setenv("REDIS_DB", "invalid")
	os.Setenv("REDIS_MIN_IDLE_CONNECTIONS", "not-a-number")
	os.Setenv("REDIS_POOL_SIZE", "abc")
	os.Setenv("REDIS_POOL_TIMEOUT", "xyz")

	defer func() {
		os.Unsetenv("APP_ENV")
		os.Unsetenv("REDIS_HOST")
		os.Unsetenv("REDIS_PORT")
		os.Unsetenv("REDIS_DB")
		os.Unsetenv("REDIS_MIN_IDLE_CONNECTIONS")
		os.Unsetenv("REDIS_POOL_SIZE")
		os.Unsetenv("REDIS_POOL_TIMEOUT")
	}()

	cfg := &config.Config{}

	client := NewRedisClient(cfg)

	if client == nil {
		t.Error("Expected non-nil Redis client")
	}

	// Client should still be created even with invalid values
	// utils.ConvertStringToInteger returns 0 for invalid strings
	opts := client.Options()

	// Just verify client was created successfully
	if opts.Addr != "localhost:6379" {
		t.Errorf("Expected address 'localhost:6379', got %s", opts.Addr)
	}
}
