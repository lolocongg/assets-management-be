package error_middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		if appErr, ok := err.(*AppError); ok {
			response := gin.H{
				"http_status": appErr.HTTPStatus,
				"code":        appErr.Code,
				"message":     appErr.Message,
			}
			if appErr.Details != nil {
				response["details"] = appErr.Details
			}
			c.JSON(appErr.HTTPStatus, response)

		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"http_status": http.StatusInternalServerError,
				"code":        CodeInternal,
				"message":     "Lỗi máy chủ",
			})
		}
		c.Abort()
	}
}
