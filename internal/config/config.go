package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
	Redis       RedisConfig
	JWT         JWTConfig
	Email       EmailConfig
	SMS         SMSConfig
	Swagger     SwaggerConfig
}

type ServerConfig struct {
	Port int
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type JWTConfig struct {
	Secret         string
	AccessExpiry   int // minutes
	RefreshExpiry  int // days
	Issuer         string
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

type SMSConfig struct {
	Provider string
	APIKey   string
	APISecret string
}

type SwaggerConfig struct {
	Enabled   bool
	Whitelist []string
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("/etc/openauth")

	// Set defaults
	viper.SetDefault("environment", "development")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("jwt.access_expiry", 15)
	viper.SetDefault("jwt.refresh_expiry", 7)
	viper.SetDefault("jwt.issuer", "openauth")
	viper.SetDefault("swagger.enabled", true)
	viper.SetDefault("swagger.whitelist", []string{})

	// Environment variables
	viper.SetEnvPrefix("OPENAUTH")
	viper.AutomaticEnv()
	
	// Bind environment variables explicitly
	viper.BindEnv("database.host", "OPENAUTH_DATABASE_HOST")
	viper.BindEnv("database.port", "OPENAUTH_DATABASE_PORT")
	viper.BindEnv("database.user", "OPENAUTH_DATABASE_USER")
	viper.BindEnv("database.password", "OPENAUTH_DATABASE_PASSWORD")
	viper.BindEnv("database.dbname", "OPENAUTH_DATABASE_DBNAME")
	viper.BindEnv("redis.host", "OPENAUTH_REDIS_HOST")
	viper.BindEnv("redis.port", "OPENAUTH_REDIS_PORT")
	viper.BindEnv("jwt.secret", "OPENAUTH_JWT_SECRET")
	viper.BindEnv("environment", "OPENAUTH_ENVIRONMENT")
	viper.BindEnv("swagger.enabled", "OPENAUTH_SWAGGER_ENABLED")
	viper.BindEnv("swagger.whitelist", "OPENAUTH_SWAGGER_WHITELIST")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	cfg := &Config{
		Environment: viper.GetString("environment"),
		Server: ServerConfig{
			Port: viper.GetInt("server.port"),
			Host: viper.GetString("server.host"),
		},
		Database: DatabaseConfig{
			Host:     getEnvOrViper("database.host", "localhost"),
			Port:     viper.GetInt("database.port"),
			User:     getEnvOrViper("database.user", "postgres"),
			Password: getEnvOrViper("database.password", "postgres"),
			DBName:   getEnvOrViper("database.dbname", "openauth"),
			SSLMode:  viper.GetString("database.ssl_mode"),
		},
		Redis: RedisConfig{
			Host:     getEnvOrViper("redis.host", "localhost"),
			Port:     viper.GetInt("redis.port"),
			Password: getEnvOrViper("redis.password", ""),
			DB:       viper.GetInt("redis.db"),
		},
		JWT: JWTConfig{
			Secret:        getEnvOrViper("jwt.secret", "change-me-in-production"),
			AccessExpiry:  viper.GetInt("jwt.access_expiry"),
			RefreshExpiry: viper.GetInt("jwt.refresh_expiry"),
			Issuer:        viper.GetString("jwt.issuer"),
		},
		Email: EmailConfig{
			SMTPHost:     getEnvOrViper("email.smtp_host", ""),
			SMTPPort:     viper.GetInt("email.smtp_port"),
			SMTPUser:     getEnvOrViper("email.smtp_user", ""),
			SMTPPassword: getEnvOrViper("email.smtp_password", ""),
			FromEmail:    getEnvOrViper("email.from_email", ""),
			FromName:     getEnvOrViper("email.from_name", "OpenAuth"),
		},
		SMS: SMSConfig{
			Provider:  getEnvOrViper("sms.provider", ""),
			APIKey:    getEnvOrViper("sms.api_key", ""),
			APISecret: getEnvOrViper("sms.api_secret", ""),
		},
		Swagger: SwaggerConfig{
			Enabled:   viper.GetBool("swagger.enabled"),
			Whitelist: getSwaggerWhitelist(),
		},
	}

	// Validate required fields
	if cfg.JWT.Secret == "change-me-in-production" && cfg.Environment == "production" {
		return nil, fmt.Errorf("JWT secret must be set in production")
	}

	return cfg, nil
}

func getEnvOrViper(key, defaultValue string) string {
	// Viper automatically converts "database.host" to "OPENAUTH_DATABASE_HOST"
	// when AutomaticEnv() is enabled with SetEnvPrefix("OPENAUTH")
	if value := viper.GetString(key); value != "" {
		return value
	}
	return defaultValue
}

// getSwaggerWhitelist parses Swagger whitelist from config or environment variable
func getSwaggerWhitelist() []string {
	// Check environment variable first (comma-separated)
	if envWhitelist := viper.GetString("swagger.whitelist"); envWhitelist != "" {
		// Split by comma and trim spaces
		whitelist := strings.Split(envWhitelist, ",")
		result := make([]string, 0, len(whitelist))
		for _, ip := range whitelist {
			ip = strings.TrimSpace(ip)
			if ip != "" {
				result = append(result, ip)
			}
		}
		if len(result) > 0 {
			return result
		}
	}

	// Check config file (YAML array)
	if whitelist := viper.GetStringSlice("swagger.whitelist"); len(whitelist) > 0 {
		return whitelist
	}

	return []string{}
}
