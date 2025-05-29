package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort  string `mapstructure:"SERVER_PORT"`
	Environment string `mapstructure:"ENV"`
	GinMode     string `mapstructure:"GIN_MODE"`

	// Database
	DatabaseURL string `mapstructure:"DATABASE_URL"`

	// JWT
	JWTSecret              string `mapstructure:"JWT_SECRET"`
	JWTExpiration          string `mapstructure:"JWT_EXPIRATION"`
	RefreshTokenExpiration string `mapstructure:"REFRESH_TOKEN_EXPIRATION"`

	// WhatsApp
	WhatsAppAPIBaseURL string `mapstructure:"WHATSAPP_API_BASE_URL"`
	WhatsAppAPIKey     string `mapstructure:"WHATSAPP_API_KEY"`
	WhatsAppSessionID  string `mapstructure:"WHATSAPP_SESSION_ID"`	

	// Admin
	AdminToken string `mapstructure:"ADMIN_TOKEN"`
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}
