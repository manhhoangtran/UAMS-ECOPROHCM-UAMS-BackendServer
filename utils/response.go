package utils

import "github.com/gin-gonic/gin"

// Http error response payload
type ErrorResponse struct {
	StatusCode int
	Msg        string
	ErrorMsg   string
}

func ResponseJson(c *gin.Context, statusCode int, data interface{}) {
	c.IndentedJSON(statusCode, data)
}
