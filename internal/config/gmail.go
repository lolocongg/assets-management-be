package config

import "os"

type GmailConfig struct {
	Email    string
	Password string
	Host     string
	Port     string
}

func LoadGmailConfig() GmailConfig {
	return GmailConfig{
		Email:    os.Getenv("GMAIL_EMAIL"),
		Password: os.Getenv("GMAIL_APP_PASSWORD"),
		Host:     os.Getenv("GMAIL_HOST"),
		Port:     os.Getenv("GMAIL_PORT"),
	}
}
