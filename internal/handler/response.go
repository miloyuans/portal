package handler

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Success bool      `json:"success"`
	Data    any       `json:"data,omitempty"`
	Error   *APIError `json:"error,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func JSONSuccess(c *gin.Context, status int, data any) {
	c.JSON(status, APIResponse{
		Success: true,
		Data:    data,
	})
}

func JSONError(c *gin.Context, status int, code, message string, details any) {
	c.JSON(status, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}
