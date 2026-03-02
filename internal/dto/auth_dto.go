package dto

import (
	"github.com/golang-jwt/jwt/v5"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required" label:"Tên đăng nhập"`
	Password string `json:"password" binding:"required,min=8" label:"Mật khẩu"`
	Role     string `json:"role" binding:"required" label:"Vai trò"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required" label:"Tên đăng nhập"`
	Password string `json:"password" binding:"required" label:"Mật khẩu"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type AuthUser struct {
	ID   int
	Role string
}

type AuthClaims struct {
	ID   int64  `json:"sub"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}
