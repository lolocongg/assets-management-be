package validator

import (
	"fmt"
	"mime/multipart"
	"reflect"
	"strings"
	// "time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func HandleValidationError(err error, dto any) map[string]any {
	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}

	errors := make(map[string]any)

	t := reflect.TypeOf(dto)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for _, fe := range ve {
		fieldKey := fe.Field()
		fieldLabel := fieldKey

		if t.Kind() == reflect.Struct {
			if sf, ok := t.FieldByName(fe.StructField()); ok {
				if jsonTag := sf.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
					fieldKey = strings.Split(jsonTag, ",")[0]
				} else if formTag := sf.Tag.Get("form"); formTag != "" {
					fieldKey = strings.Split(formTag, ",")[0]
				}
				if label := sf.Tag.Get("label"); label != "" {
					fieldLabel = label
				}
			}
		}

		switch fe.Tag() {
		case "required":
			errors[fieldKey] = fieldLabel + " là bắt buộc"
		case "min":
			errors[fieldKey] = fieldLabel + " phải có tối thiểu " + fe.Param() + " ký tự"
		case "max":
			errors[fieldKey] = fieldLabel + " phải có tối đa " + fe.Param() + " ký tự"
		case "oneof":
			errors[fieldKey] = fieldLabel + " phải là một trong các giá trị: " + fe.Param()
		case "gte":
			errors[fieldKey] = fieldLabel + " phải lớn hơn hoặc bằng " + fe.Param()
		case "lte":
			errors[fieldKey] = fieldLabel + " phải nhỏ hơn hoặc bằng " + fe.Param()
		case "gt":
			errors[fieldKey] = fieldLabel + " phải lớn hơn " + fe.Param()
		case "lt":
			errors[fieldKey] = fieldLabel + " phải nhỏ hơn " + fe.Param()
		case "gtefield":
			errors[fieldKey] = fieldLabel + " không được trước " + fe.Param()
		case "images":
			errors[fieldKey] = fieldLabel + " không đúng format"
		default:
			errors[fieldKey] = fieldLabel + " không hợp lệ"
		}
	}

	return errors
}

func IsValidImageMime(mime string) bool {
	allowed := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/jpg":  true,
	}
	return allowed[mime]
}

func RegisterValidators() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return fmt.Errorf("validator engine not supported")
	}

	// Register image validator
	v.RegisterValidation("images", func(fl validator.FieldLevel) bool {
		files, ok := fl.Field().Interface().([]*multipart.FileHeader)
		if !ok {
			return false
		}

		if len(files) == 0 {
			return true
		}

		for _, file := range files {
			if file.Size > 5<<20 { // 5MB
				return false
			}

			mime := file.Header.Get("Content-Type")
			if !IsValidImageMime(mime) {
				return false
			}
		}
		return true
	})

	// Register tag name function for label, json, and form tags
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		if label := fld.Tag.Get("label"); label != "" {
			return label
		}

		if jsonTag := fld.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
			return strings.Split(jsonTag, ",")[0]
		}

		if formTag := fld.Tag.Get("form"); formTag != "" {
			return formTag
		}

		return fld.Name
	})
	return nil
}
