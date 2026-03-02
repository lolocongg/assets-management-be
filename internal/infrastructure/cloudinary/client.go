package cloudinary

import (
	"fmt"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/davidcm146/assets-management-be.git/internal/config"
)

func NewCloudinary(cfg *config.CloudinaryConfig) (*cloudinary.Cloudinary, error) {
	cld, err := cloudinary.NewFromParams(
		cfg.CloudName,
		cfg.APIKey,
		cfg.APISecret,
	)
	if err != nil {
		return nil, fmt.Errorf("cloudinary init failed: %w", err)
	}
	return cld, nil
}
