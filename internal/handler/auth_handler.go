package handler

import (
	"net/http"

	"github.com/davidcm146/assets-management-be.git/internal/dto"
	"github.com/davidcm146/assets-management-be.git/internal/error_middleware"
	"github.com/davidcm146/assets-management-be.git/internal/model"
	"github.com/davidcm146/assets-management-be.git/internal/repository"
	"github.com/davidcm146/assets-management-be.git/internal/service"
	"github.com/davidcm146/assets-management-be.git/internal/validator"
	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	LoginHandler(c *gin.Context)
	RegisterHandler(c *gin.Context)
	GetMe(c *gin.Context)
}

type authHandler struct {
	authService service.AuthService
	userRepo    repository.UserRepository
}

func NewAuthHandler(authService service.AuthService, userRepo repository.UserRepository) AuthHandler {
	return &authHandler{authService: authService, userRepo: userRepo}
}

func (h *authHandler) LoginHandler(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(error_middleware.
			NewUnprocessableEntity("Dữ liệu không hợp lệ").
			WithDetails(validator.HandleValidationError(err, &req)),
		)
		return
	}
	token, err := h.authService.LoginService(c.Request.Context(), req.Username, req.Password)

	if err != nil {
		if _, ok := err.(*error_middleware.AppError); ok {
			c.Error(err)
			return
		}

		c.Error(error_middleware.NewInternal("Lỗi hệ thống"))
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponse{Token: token})
}

func (h *authHandler) RegisterHandler(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(
			error_middleware.
				NewUnprocessableEntity("Dữ liệu không hợp lệ").
				WithDetails(validator.HandleValidationError(err, &req)),
		)
		return
	}

	user := &model.User{
		Username: req.Username,
		Password: req.Password,
		Role:     model.ParseRole(req.Role),
	}

	if err := h.authService.RegisterService(c.Request.Context(), user); err != nil {
		if _, ok := err.(*error_middleware.AppError); ok {
			c.Error(err)
			return
		}

		c.Error(error_middleware.NewInternal("Đăng ký không thành công"))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Tạo tài khoản thành công"})
}

func (h *authHandler) GetMe(c *gin.Context) {
	userRaw, exists := c.Get("user")
	if !exists {
		c.Error(error_middleware.NewUnauthorized("Không tìm thấy thông tin người dùng"))
		return
	}

	authUser, ok := userRaw.(*dto.AuthUser)
	if !ok {
		c.Error(error_middleware.NewUnauthorized("Thông tin người dùng không hợp lệ"))
		return
	}
	user, err := h.userRepo.GetByID(c.Request.Context(), authUser.ID)
	if err != nil {
		c.Error(error_middleware.NewNotFound("Người dùng không tồn tại"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
	})
}
