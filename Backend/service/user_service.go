package service

import (
	"context"
	"time"

	"fullstack-go-grpc/database/internals/models"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// UserService provides business logic for user operations.
type UserService interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	ListUsers(ctx context.Context) ([]models.User, error)
}

type userServiceImpl struct {
	db *bun.DB
}

// NewUserService creates a new UserService.
func NewUserService(db *bun.DB) UserService {
	return &userServiceImpl{db: db}
}

func (s *userServiceImpl) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	_, err := s.db.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userServiceImpl) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := s.db.NewSelect().Model(&user).Where("unique_id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *userServiceImpl) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	user.UpdatedAt = time.Now()
	_, err := s.db.NewUpdate().Model(user).Where("id = ?", user.ID).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userServiceImpl) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*models.User)(nil)).Where("unique_id = ?", id).Exec(ctx)
	return err
}

func (s *userServiceImpl) ListUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := s.db.NewSelect().Model(&users).Order("created_at DESC").Scan(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}
