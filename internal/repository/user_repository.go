package repository

import (
	"context"
	"time"

	"github.com/davidcm146/assets-management-be.git/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(ctx context.Context, u *model.User) error
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByID(ctx context.Context, id int) (*model.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *model.User) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	_, err := r.db.Exec(ctx,
		"INSERT INTO users (username, role, password, created_at) VALUES ($1, $2, $3, $4)",
		u.Username, u.Role, string(hash), time.Now().UTC(),
	)
	return err
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	row := r.db.QueryRow(ctx,
		"SELECT id, username, role, password FROM users WHERE username=$1",
		username,
	)
	var user model.User
	err := row.Scan(&user.ID, &user.Username, &user.Role, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	row := r.db.QueryRow(ctx,
		"SELECT id, username, role FROM users WHERE id=$1",
		id,
	)
	var user model.User
	err := row.Scan(&user.ID, &user.Username, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
