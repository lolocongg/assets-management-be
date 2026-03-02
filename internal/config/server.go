package config

import "os"

type ServerConfig struct {
	Port string
}

func LoadServerConfig() ServerConfig {
	return ServerConfig{
		Port: os.Getenv("PORT"),
	}
}
