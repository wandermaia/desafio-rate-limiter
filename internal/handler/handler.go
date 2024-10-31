package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// handler de teste
func TestHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "request successful"})
}
