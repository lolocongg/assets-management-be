package service

import (
	"context"
	"os"
	"time"

	"github.com/davidcm146/assets-management-be.git/internal/error_middleware"
	"github.com/davidcm146/assets-management-be.git/internal/model"
	"github.com/davidcm146/assets-management-be.git/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	LoginService(ctx context.Context, username, password string) (string, error)
	RegisterService(ctx context.Context, u *model.User) error
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) LoginService(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", error_middleware.NewInternal("Lỗi hệ thống")
	}
	if user == nil {
		return "", error_middleware.NewUnauthorized("Tài khoản hoặc mật khẩu không đúng")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", error_middleware.NewUnauthorized("Tài khoản hoặc mật khẩu không đúng")
	}

	claims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"role":     user.Role.String(),
		"exp":      jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", error_middleware.NewInternal("Không thể tạo token đăng nhập")
	}

	return tokenString, nil
}

func (s *authService) RegisterService(ctx context.Context, u *model.User) error {
	err := s.userRepo.Create(ctx, u)
	if err != nil {
		return error_middleware.NewInternal("Không thể tạo tài khoản")
	}

	return nil
}
