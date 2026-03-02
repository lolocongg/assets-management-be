package config

import (
	"github.com/joho/godotenv"
)

type Config struct {
	Server     ServerConfig
	DB         DBConfig
	Cloudinary CloudinaryConfig
	Postmark   PostmarkConfig
	Gmail      GmailConfig
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Server:     LoadServerConfig(),
		DB:         LoadDatabaseConfig(),
		Cloudinary: LoadCloudinaryConfig(),
		Postmark:   LoadPostmarkConfig(),
		Gmail:      LoadGmailConfig(),
	}

	return cfg, nil
}
