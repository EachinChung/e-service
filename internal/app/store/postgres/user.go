package postgres

import (
	"context"

	"gorm.io/gorm"

	"github.com/eachinchung/component-base/db/options"
	"github.com/eachinchung/errors"

	"github.com/eachinchung/e-service/internal/app/store"
	"github.com/eachinchung/e-service/internal/app/store/model"
)

type user struct{}

func newUser() *user {
	return &user{}
}

var _ store.UserStore = &user{}

func (u user) Create(ctx context.Context, db *gorm.DB, user *model.Users) error {
	if err := db.Create(user).Error; err != nil {
		return errors.Wrap(err, "failed to create user")
	}
	return nil
}

func (u user) Delete(ctx context.Context, db *gorm.DB, user *model.Users, opts ...options.Opt) error {
	o := &options.Option{Unscoped: false}

	for _, opt := range opts {
		opt(o)
	}

	if o.Unscoped {
		db = db.Unscoped()
	}

	if err := db.Delete(user).Error; err != nil {
		return errors.Wrap(err, "failed to delete user")
	}
	return nil
}

func (u user) Get(ctx context.Context, db *gorm.DB, key any, opts ...options.Opt) (*model.Users, error) {
	o := &options.Option{
		Unscoped: false,
		Where: options.Where{
			Query: "eid = ?",
			Args:  []any{key},
		},
	}

	for _, opt := range opts {
		opt(o)
	}

	if o.Unscoped {
		db = db.Unscoped()
	}

	var user model.Users
	if err := db.Where(o.Where.Query, o.Where.Args...).First(&user).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get user")
	}
	return &user, nil
}

func (u user) Update(ctx context.Context, db *gorm.DB, user *model.Users) error {
	if err := db.Save(user).Error; err != nil {
		return errors.Wrap(err, "failed to update user")
	}
	return nil
}
