package validator

import (
	"errors"
	"mentorship-app-backend/components/errorpackage"
	"regexp"
)

func ValidateKey(key string) error {
	if key == "" {
		return errorpackage.ErrNoSuchKey
	}
	return nil
}

func ValidateEmail(email string) error {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func ValidateName(name string) error {
	if len(name) == 0 {
		return errors.New("name is required")
	}
	if len(name) < 2 {
		return errors.New("name must be at least 2 characters long")
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	const lowercaseRegex = `[a-z]`
	if !regexp.MustCompile(lowercaseRegex).MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}

	const uppercaseRegex = `[A-Z]`
	if !regexp.MustCompile(uppercaseRegex).MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}

	const digitRegex = `\d`
	if !regexp.MustCompile(digitRegex).MatchString(password) {
		return errors.New("password must contain at least one number")
	}

	return nil
}

func ValidateRole(role string) error {
	if role != "mentor" && role != "mentee" {
		return errors.New("invalid role; must be either 'mentor' or 'mentee'")
	}
	return nil
}

func ValidateFields(name, email, password, role string) error {
	if err := ValidateName(name); err != nil {
		return err
	}
	if err := ValidateEmail(email); err != nil {
		return err
	}
	if err := ValidatePassword(password); err != nil {
		return err
	}
	if err := ValidateRole(role); err != nil {
		return err
	}
	return nil
}
