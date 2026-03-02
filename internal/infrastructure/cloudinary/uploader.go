package cloudinary

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryUploader struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryUploader(cld *cloudinary.Cloudinary) *CloudinaryUploader {
	return &CloudinaryUploader{cld: cld}
}

func (u *CloudinaryUploader) Upload(ctx context.Context, file *multipart.FileHeader) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	resp, err := u.cld.Upload.Upload(ctx, f,
		uploader.UploadParams{
			Folder: "loan_slips",
		},
	)
	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}
