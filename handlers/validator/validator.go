package validator

import (
	"mentorship-app-backend/handlers/errorpackage"
)

func ValidateKey(key string) error {
	if key == "" {
		return errorpackage.ErrKeyNotFound
	}

	return nil
}

func IsValidMimeType(mimeType string) bool {
	validMimeTypes := []string{"image/jpeg", "image/png", "application/pdf", "text/plain"}
	for _, valid := range validMimeTypes {
		if mimeType == valid {
			return true
		}
	}
	return false
}
