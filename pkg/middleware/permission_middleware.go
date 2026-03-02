package middleware

import (
	"net/http"

	"github.com/davidcm146/assets-management-be.git/internal/error_middleware"
	"github.com/davidcm146/assets-management-be.git/internal/model"
	"github.com/gin-gonic/gin"
)

func PermissionMiddleware(requiredRole []model.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.MustGet("role")

		isAllowed := false
		for _, r := range requiredRole {
			if role == r.String() {
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, &error_middleware.AppError{
				HTTPStatus: http.StatusForbidden,
				Code:       error_middleware.CodeForbidden,
				Message:    "Bạn không có quyền truy cập tài nguyên này",
			})
			return
		}
		c.Next()
	}
}
