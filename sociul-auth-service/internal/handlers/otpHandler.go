package handlers

import (
	"net/http"
	"sociul-auth-service/internal/models"
	"sociul-auth-service/internal/sentinel"
	"sociul-auth-service/internal/service"

	"github.com/gin-gonic/gin"
)

type OtpHandler struct {
	service *service.OtpService
}

func NewOtpHandler(service *service.OtpService) *OtpHandler {
	return &OtpHandler{service: service}
}

func (h *OtpHandler) SendEmailOtp(c *gin.Context) {
	// Extract request body
	var request models.SendOtpRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Call service layer
	h.service.SendEmailOtp(c.Request.Context(), request)

	// Always send this, even if email does not exist, to prevent email enumeration
	c.JSON(http.StatusOK, gin.H{"message": "if an account exists, an otp was sent"})
}

func (h *OtpHandler) VerifyEmailOtp(c *gin.Context) {
	// Extract request body
	var request models.VerifyOtpRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Call service layer
	err := h.service.VerifyEmailOtp(c.Request.Context(), request)
	if err != nil {
		sentinel.WriteResp(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "email verified successfully"})
}

func (h *OtpHandler) SendResetOtp(c *gin.Context) {
	// Extract request body
	var request models.SendOtpRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Call service layer
	h.service.SendResetOtp(c.Request.Context(), request)

	// Always send this, even if email does not exist, to prevent email enumeration
	c.JSON(http.StatusOK, gin.H{"message": "if an account exists, an otp was sent"})
}

func (h *OtpHandler) VerifyResetOtp(c *gin.Context) {
	// Extract request body
	var request models.VerifyOtpRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Call service layer
	token, err := h.service.VerifyResetOtp(c.Request.Context(), request)
	if err != nil {
		sentinel.WriteResp(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"reset-token": token})
}
