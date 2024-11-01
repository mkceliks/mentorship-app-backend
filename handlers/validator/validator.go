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
