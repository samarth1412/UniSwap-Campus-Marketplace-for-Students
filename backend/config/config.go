package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config stores runtime configuration loaded from environment variables.
type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	JWTSecret  string
}

func Load() (*Config, error) {
	fileValues, err := loadResourceEnv(filepath.Join("resources", "app.env"))
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Port:       getConfigValue(fileValues, "PORT", "8080"),
		DBHost:     getConfigValue(fileValues, "DB_HOST", "localhost"),
		DBPort:     getConfigValue(fileValues, "DB_PORT", "5432"),
		DBUser:     getConfigValue(fileValues, "DB_USER", "postgres"),
		DBPassword: getConfigValue(fileValues, "DB_PASSWORD", ""),
		DBName:     getConfigValue(fileValues, "DB_NAME", "uniswap"),
		DBSSLMode:  getConfigValue(fileValues, "DB_SSLMODE", "disable"),
		JWTSecret:  getConfigValue(fileValues, "JWT_SECRET", ""),
	}

	if cfg.DBPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD is required")
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func (c *Config) DatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getConfigValue(fileValues map[string]string, key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	if value, ok := fileValues[key]; ok && value != "" {
		return value
	}
	return fallback
}

func loadResourceEnv(path string) (map[string]string, error) {
	values := make(map[string]string)

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return values, nil
		}
		return nil, fmt.Errorf("open config file %s: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, "\"")
		values[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read config file %s: %w", path, err)
	}

	return values, nil
}
