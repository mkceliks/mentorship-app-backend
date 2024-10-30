package validator

import errorPackage "mentorship-app-backend/handlers/s3/errorpackage"

func ValidateKey(key string) error {
	if key == "" {
		return errorPackage.ErrKeyNotFound
	}

	return nil
}
