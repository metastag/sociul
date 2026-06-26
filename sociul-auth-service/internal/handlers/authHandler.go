package handlers

import (
	"net/http"
	"strings"

	"sociul-auth-service/internal/models"
	"sociul-auth-service/internal/sentinel"
	"sociul-auth-service/internal/service"
	"sociul-auth-service/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Middleware to verify jwt token and write user id to context
func (h *AuthHandler) VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
			return
		}

		// Check if header is in correct format
		if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Extract token string from header
		tokenString := strings.Split(authHeader, "Bearer ")[1]

		// Send token to service layer to verify and extract user id
		id, err := h.service.ValidateJwt(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Create a new context with user id and set it as requext context
		ctxWithUser := utils.NewContext(c.Request.Context(), id)
		c.Request = c.Request.WithContext(ctxWithUser)
		c.Next()
	}
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	// Extract request body
	var request models.SignUpRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Call service layer
	refreshToken, accessToken, err := h.service.SignUp(c.Request.Context(), request)
	if err != nil {
		sentinel.WriteResp(c, err)
		return
	}

	// Return Refresh and Access token
	c.JSON(http.StatusCreated, gin.H{"refreshToken": refreshToken, "accessToken": accessToken})
}

func (h *AuthHandler) Login(c *gin.Context) {
	// Extract request body
	var request models.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Client Ip is used to track failed login attempts per ip
	clientIp := c.ClientIP()

	// Call service layer
	refreshToken, accessToken, err := h.service.Login(c.Request.Context(), request, clientIp)
	if err != nil {
		sentinel.WriteResp(c, err)
		return
	}

	// Return Refresh and Access token
	c.JSON(http.StatusOK, gin.H{"refreshToken": refreshToken, "accessToken": accessToken})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	// Extract request body
	var request models.RefreshRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Call service layer
	accessToken, err := h.service.Refresh(c.Request.Context(), request)
	if err != nil {
		sentinel.WriteResp(c, err)
		return
	}

	// Return new access token
	c.JSON(http.StatusOK, gin.H{"accessToken": accessToken})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Extract request body
	var request models.RefreshRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Call service layer
	err := h.service.Logout(c.Request.Context(), request)
	if err != nil {
		sentinel.WriteResp(c, err)
		return
	}

	// Return new access token
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})

}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	// Extract request body
	var request models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Call service layer
	err := h.service.ResetPassword(c.Request.Context(), request)
	if err != nil {
		sentinel.WriteResp(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password successfully reset!"})
}

func (h *AuthHandler) Delete(c *gin.Context) {
	// Call service layer
	err := h.service.DeleteUser(c.Request.Context())
	if err != nil {
		sentinel.WriteResp(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
