package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/davidcm146/assets-management-be.git/internal/dto"
	"github.com/davidcm146/assets-management-be.git/internal/error_middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.Error(&error_middleware.AppError{
				HTTPStatus: http.StatusUnauthorized,
				Code:       error_middleware.CodeUnauthorized,
				Message:    "Yêu cầu xác thực",
			})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		if tokenStr == "" {
			c.Error(&error_middleware.AppError{
				HTTPStatus: http.StatusUnauthorized,
				Code:       error_middleware.CodeUnauthorized,
				Message:    "Yêu cầu xác thực",
			})
			c.Abort()
			return
		}

		claims := &dto.AuthClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.Error(&error_middleware.AppError{
				HTTPStatus: http.StatusUnauthorized,
				Code:       error_middleware.CodeUnauthorized,
				Message:    "Token không hợp lệ",
			})
			c.Abort()
			return
		}
		c.Set("role", claims.Role)
		c.Set("user", &dto.AuthUser{
			ID:   int(claims.ID),
			Role: claims.Role,
		})

		c.Next()
	}
}
