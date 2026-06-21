package sentinel

import "errors"

var (
	// Auth errors
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidToken = errors.New("invalid token")

	// Otp/token cache errors
	ErrCacheMiss = errors.New("cache miss")
	ErrWrongOtp  = errors.New("wrong otp")

	// Username, password, email errors
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidUsername       = errors.New("username length not in 5-255 characters")
	ErrUsernameSpecial       = errors.New("Username contains special characters")
	ErrUsernameAlreadyExists = errors.New("A user with this username already exists")
	ErrInvalidPassword       = errors.New("password length not in 8-70 characters")
	ErrWrongPassword         = errors.New("wrong password")
	ErrInvalidEmail          = errors.New("email length not in 7-255 characters")
	ErrEmailPattern          = errors.New("invalid email format")

	// Generic internal server error
	ErrInternal = errors.New("internal server error")
)
