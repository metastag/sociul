package utils

import (
	"regexp"
	"strings"

	"sociul-auth-service/internal/models"
	"sociul-auth-service/internal/sentinel"
)

// Standard normalization for user input
func SignUpNormalize(user models.SignUpRequest) models.SignUpRequest {
	// Remove trailing whitespace
	user.Username = strings.TrimSpace(user.Username)
	user.Password = strings.TrimSpace(user.Password)
	user.Email = strings.TrimSpace(user.Email)

	// Username and email should always be lowercase
	user.Username = strings.ToLower(user.Username)
	user.Email = strings.ToLower(user.Email)

	return user
}

func LoginNormalize(user models.LoginRequest) models.LoginRequest {
	// Remove trailing whitespcae
	user.Username = strings.TrimSpace(user.Username)
	user.Password = strings.TrimSpace(user.Password)

	// Username should always be lowercase
	user.Username = strings.ToLower(user.Username)

	return user
}

// Validation for user input
func SignUpValidate(user models.SignUpRequest) error {
	// Length check
	if len(user.Username) < 5 || len(user.Username) > 254 {
		return sentinel.ErrInvalidUsername
	}

	if len(user.Password) < 8 || len(user.Password) > 70 {
		return sentinel.ErrInvalidPassword
	}

	if len(user.Email) < 7 || len(user.Email) > 254 {
		return sentinel.ErrInvalidEmail
	}

	// Email pattern validation
	match, err := regexp.MatchString(`^[^@]+@[^@]+\.[^@]+$`, user.Email)
	if err != nil || !match {
		return sentinel.ErrEmailPattern
	}

	// Dont allow special characters in username, except hyphen and underscore
	match, err = regexp.MatchString("^[a-zA-Z0-9_-]+$", user.Username)
	if err != nil || !match {
		return sentinel.ErrUsernameSpecial
	}

	return nil
}

func LoginValidate(user models.LoginRequest) error {
	// Length check
	if len(user.Username) < 5 || len(user.Username) > 254 {
		return sentinel.ErrInvalidUsername
	}

	if len(user.Password) < 8 || len(user.Password) > 70 {
		return sentinel.ErrInvalidPassword
	}

	// Dont allow special characters in username, except hyphen and underscore
	match, err := regexp.MatchString("^[a-zA-Z0-9_-]+$", user.Username)
	if err != nil || !match {
		return sentinel.ErrUsernameSpecial
	}

	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 || len(password) > 70 {
		return sentinel.ErrInvalidPassword
	}
	return nil
}
