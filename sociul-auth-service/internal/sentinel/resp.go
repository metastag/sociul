package sentinel

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Helper function to map sentinel errors to appropriate http responses
func WriteResp(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	case errors.Is(err, ErrInvalidToken):
		c.JSON(http.StatusBadRequest, gin.H{"error": "token expired, please log in again"})
	case errors.Is(err, ErrCacheMiss):
		c.JSON(http.StatusBadRequest, gin.H{"error": "otp or token expired, please try again"})
	case errors.Is(err, ErrWrongOtp):
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid otp or token"})
	case errors.Is(err, ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
	case errors.Is(err, ErrInvalidUsername):
		c.JSON(http.StatusBadRequest, gin.H{"error": "username must be between 5 and 255 characters long"})
	case errors.Is(err, ErrUsernameSpecial):
		c.JSON(http.StatusBadRequest, gin.H{"error": "username must not contain special characters"})
	case errors.Is(err, ErrUsernameAlreadyExists):
		c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
	case errors.Is(err, ErrInvalidPassword):
		c.JSON(http.StatusBadRequest, gin.H{"error": "password must be between 8 and 70 characters long"})
	case errors.Is(err, ErrWrongPassword):
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong password"})
	case errors.Is(err, ErrInvalidEmail):
		c.JSON(http.StatusBadRequest, gin.H{"error": "email must be between 7 and 255 characters long"})
	case errors.Is(err, ErrEmailPattern):
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format, try again"})
	case errors.Is(err, ErrInternal):
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	default:
		log.Println("Unexpected error received in writeResp() - ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
