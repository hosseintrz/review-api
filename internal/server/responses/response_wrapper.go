package responses

import (
	"github.com/gin-gonic/gin"
)

type Error struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type Data struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func AbortResponse(c *gin.Context, statusCode int, data interface{}) {
	c.AbortWithStatusJSON(statusCode, data)
}

func Response(c *gin.Context, statusCode int, data interface{}) {
	// nolint // context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	// nolint // context.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// nolint // context.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization")
	c.JSON(statusCode, data)
}

func MessageResponse(c *gin.Context, statusCode int, message string) {
	Response(c, statusCode, Data{
		Code:    statusCode,
		Message: message,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, err error) {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	c.AbortWithStatusJSON(statusCode, gin.H{"error": errorMsg})
}
