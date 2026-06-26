package sentinel

import "net/http"

type AppError struct {
	Code    int
	Message string
	Wrapped error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Wrapped
}

func NewAppError(code int, message string, wrapped error) *AppError {
	return &AppError{Code: code, Message: message, Wrapped: wrapped}
}

var (
	// Auth errors
	ErrUnauthorized = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "unauthorized",
	}
	ErrInvalidToken = &AppError{
		Code:    http.StatusBadRequest,
		Message: "token expired, please log in again",
	}
	ErrRateLimit = &AppError{
		Code:    http.StatusBadRequest,
		Message: "too many login attempts, try again later",
	}

	// Otp/token cache errors
	ErrCacheMiss = &AppError{
		Code:    http.StatusBadRequest,
		Message: "otp or token expired, please try again",
	}
	ErrWrongOtp = &AppError{
		Code:    http.StatusBadRequest,
		Message: "invalid otp or token",
	}

	// Username, password, email errors
	ErrUserNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "invalid username or password",
	}
	ErrInvalidUsername = &AppError{
		Code:    http.StatusBadRequest,
		Message: "username must be between 5 and 255 characters long",
	}
	ErrUsernameSpecial = &AppError{
		Code:    http.StatusBadRequest,
		Message: "username must not contain special characters",
	}
	ErrUsernameAlreadyExists = &AppError{
		Code:    http.StatusConflict,
		Message: "username already exists",
	}
	ErrInvalidPassword = &AppError{
		Code:    http.StatusBadRequest,
		Message: "password must be between 8 and 70 characters long",
	}
	ErrWrongPassword = &AppError{
		Code:    http.StatusBadRequest,
		Message: "invalid username or password",
	}
	ErrInvalidEmail = &AppError{
		Code:    http.StatusBadRequest,
		Message: "email must be between 7 and 255 characters long",
	}
	ErrEmailPattern = &AppError{
		Code:    http.StatusBadRequest,
		Message: "invalid email format, try again",
	}

	// Generic internal server error
	ErrInternal = &AppError{
		Code:    http.StatusInternalServerError,
		Message: "internal server error",
	}
)
