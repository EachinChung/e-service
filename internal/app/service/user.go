package service

import (
	"context"
	"regexp"

	"gorm.io/gorm"

	"github.com/eachinchung/component-base/db/options"
	"github.com/eachinchung/errors"

	"github.com/eachinchung/e-service/internal/app/store"
	"github.com/eachinchung/e-service/internal/app/store/model"
	"github.com/eachinchung/e-service/internal/pkg/code"
)

// UserSrv defines functions used to handle user request.
type UserSrv interface {
	Create(ctx context.Context, user *model.Users) error
	GetByUsername(ctx context.Context, username string) (*model.Users, error)
}

type userService struct {
	store store.Store
}

var _ UserSrv = &userService{}

func newUsers(srv *service) *userService {
	return &userService{store: srv.store}
}

func (u userService) Create(ctx context.Context, user *model.Users) error {
	db := u.store.DB()

	if _, err := u.store.User().Get(ctx, db, user.Phone, options.WithQuery("phone = ?")); err == nil {
		return errors.Code(code.ErrPhoneAlreadyExist, "phone already exists")
	}

	if _, err := u.store.User().Get(ctx, db, user.Username, options.WithQuery("username = ?")); err == nil {
		return errors.Code(code.ErrUsernameAlreadyExist, "username already exists")
	}

	if err := u.store.User().Create(ctx, db, user); err != nil {
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key '.*'", err.Error()); match {
			return errors.Code(code.ErrUserAlreadyExist, err.Error())
		}

		return errors.Code(code.ErrDatabase, err.Error())
	}

	return nil
}

func (u userService) GetByUsername(ctx context.Context, username string) (*model.Users, error) {
	db := u.store.DB()

	user, err := u.store.User().Get(ctx, db, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Code(code.ErrUserNotExist, err.Error())
		}
		return nil, errors.Code(code.ErrDatabase, err.Error())
	}

	return user, nil
}
