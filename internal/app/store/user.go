package store

import (
	"context"

	"gorm.io/gorm"

	"github.com/eachinchung/component-base/db/options"

	"github.com/eachinchung/e-service/internal/app/store/model"
)

type UserStore interface {
	Create(ctx context.Context, db *gorm.DB, user *model.Users) error
	Delete(ctx context.Context, db *gorm.DB, user *model.Users, opts ...options.Opt) error
	Get(ctx context.Context, db *gorm.DB, key interface{}, opts ...options.Opt) (*model.Users, error)
	Update(ctx context.Context, db *gorm.DB, user *model.Users) error
}

type SuperUsersStore interface {
	Exist(ctx context.Context, db *gorm.DB, username string) (bool, error)
}
