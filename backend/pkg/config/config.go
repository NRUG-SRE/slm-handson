package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Server    ServerConfig
	NewRelic  NewRelicConfig
	Performance PerformanceConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type NewRelicConfig struct {
	APIKey  string
	AppName string
}

type PerformanceConfig struct {
	ErrorRate         float64
	ResponseTimeMin   int
	ResponseTimeMax   int
	SlowEndpointRate  float64
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Host: getEnv("HOST", "0.0.0.0"),
		},
		NewRelic: NewRelicConfig{
			APIKey:  getEnv("NEW_RELIC_API_KEY", ""),
			AppName: getEnv("NEW_RELIC_APP_NAME", "slm-handson-api"),
		},
		Performance: PerformanceConfig{
			ErrorRate:        getEnvFloat("ERROR_RATE", 0.1),
			ResponseTimeMin:  getEnvInt("RESPONSE_TIME_MIN", 50),
			ResponseTimeMax:  getEnvInt("RESPONSE_TIME_MAX", 500),
			SlowEndpointRate: getEnvFloat("SLOW_ENDPOINT_RATE", 0.2),
		},
	}
}

func (c *Config) Validate() error {
	if c.Server.Port == "" {
		log.Println("Warning: PORT not set, using default 8080")
	}
	
	if c.NewRelic.APIKey == "" {
		log.Println("Warning: NEW_RELIC_API_KEY not set, New Relic monitoring will be disabled")
	}
	
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf("Warning: Invalid integer value for %s: %s, using default %d", key, value, defaultValue)
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
		log.Printf("Warning: Invalid float value for %s: %s, using default %f", key, value, defaultValue)
	}
	return defaultValue
}