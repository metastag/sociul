package sentinel

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Helper function to map sentinel errors to appropriate http responses
func WriteResp(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		return
	}

	log.Println("Unexpected error received in WriteResp() - ", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}
