package config

import "os"

type CloudinaryConfig struct {
	CloudName string
	APIKey    string
	APISecret string
}

func LoadCloudinaryConfig() CloudinaryConfig {
	return CloudinaryConfig{
		CloudName: os.Getenv("CLOUDINARY_CLOUD_NAME"),
		APIKey:    os.Getenv("CLOUDINARY_API_KEY"),
		APISecret: os.Getenv("CLOUDINARY_API_SECRET"),
	}
}
