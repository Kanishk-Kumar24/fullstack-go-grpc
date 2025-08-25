package repo

import (
	"context"
	"fullstack-go-grpc/internals/models"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// type UserRepo interface {
// 	CreateUser(ctx context.Context, user *models.User) error
// 	GetUserByID(ctx context.Context, id uuid.UUID, user *models.User) error
// 	UpdateUser(ctx context.Context, user *models.User) error
// 	DeleteUser(ctx context.Context, id uuid.UUID) error
// 	ListUsers(ctx context.Context, user []models.User) error
// }
type UserRepo struct {
	db *bun.DB
}

func NewUserRepo(db *bun.DB) *UserRepo  {
	return &UserRepo{db: db}
}

func (r *UserRepo ) CreateUser(ctx context.Context, user *models.User) error {
	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	return err
}
func (r *UserRepo ) GetUserByID(ctx context.Context, id uuid.UUID, user *models.User) error {
	err := r.db.NewSelect().Model(&user).Where("unique_id = ?", id).Scan(ctx)
	return err
}
func (r *UserRepo ) UpdateUser(ctx context.Context, user *models.User) error {
	_, err := r.db.NewUpdate().Model(user).Where("id = ?", user.ID).Exec(ctx)
	return err
}
func (r *UserRepo ) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().Model((*models.User)(nil)).Where("unique_id = ?", id).Exec(ctx)
	return err
}
func (r *UserRepo ) ListUsers(ctx context.Context, users []models.User) error {
	err := r.db.NewSelect().Model(&users).Order("created_at DESC").Scan(ctx)
	return err
}
