package middleware

import "github.com/gin-gonic/gin"

type apiErrorResponse struct {
	Success bool         `json:"success"`
	Error   apiErrorBody `json:"error"`
}

type apiErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func abortJSON(c *gin.Context, status int, code, message string, details any) {
	c.AbortWithStatusJSON(status, apiErrorResponse{
		Success: false,
		Error: apiErrorBody{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}
