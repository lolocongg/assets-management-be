package model

import (
	"time"

	"github.com/davidcm146/assets-management-be.git/internal/utils"
)

type Role int

const (
	Admin Role = iota + 1 // Admin = 1
	IT                    // IT = 2
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username" binding:"required,gte=3,lte=10" label:"Tên đăng nhập"`
	Password  string    `json:"password" binding:"required,gte=6,lte=20" label:"Mật khẩu"`
	Role      Role      `json:"role" binding:"required,oneof=1 2" label:"Quyền hạn"`
	CreatedAt time.Time `json:"created_at"`
}

func (r Role) String() string {
	switch r {
	case Admin:
		return "admin"
	case IT:
		return "IT"
	default:
		return "Unknown"
	}
}

func ParseRole(s string) Role {
	switch s {
	case "admin":
		return Admin
	case "IT":
		return IT
	default:
		return 0
	}
}

func (r Role) MarshalJSON() ([]byte, error) {
	return utils.MarshalEnum(r)
}
