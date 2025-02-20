// config/config.go
package config

import (
    "os"
    "strconv"
)

type Config struct {
    AdditionTime       int
    SubtractionTime    int
    MultiplicationTime int
    DivisionTime       int
}

func Load() *Config {
    return &Config{
        AdditionTime:       getEnvAsInt("TIME_ADDITION_MS", 1000),
        SubtractionTime:    getEnvAsInt("TIME_SUBTRACTION_MS", 1000),
        MultiplicationTime: getEnvAsInt("TIME_MULTIPLICATIONS_MS", 1000),
        DivisionTime:       getEnvAsInt("TIME_DIVISIONS_MS", 1000),
    }
}

func getEnvAsInt(key string, defaultValue int) int {
    if value, exists := os.LookupEnv(key); exists {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}