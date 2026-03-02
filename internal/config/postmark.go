package config

import "os"

type PostmarkConfig struct {
	APIURL string
	Token  string
	From   string
}

func LoadPostmarkConfig() PostmarkConfig {
	return PostmarkConfig{
		APIURL: os.Getenv("POSTMARK_API_URL"),
		Token:  os.Getenv("POSTMARK_TOKEN"),
		From:   os.Getenv("POSTMARK_FROM"),
	}
}
