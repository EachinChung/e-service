package service

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/eachinchung/component-base/db/options"
	"github.com/eachinchung/errors"
	"github.com/eachinchung/log"
	"gorm.io/gorm"

	"github.com/eachinchung/e-service/internal/app/storage"
	"github.com/eachinchung/e-service/internal/app/store"
	"github.com/eachinchung/e-service/internal/app/store/model"
	"github.com/eachinchung/e-service/internal/pkg/code"
)

// UserSrv defines functions used to handle user request.
type UserSrv interface {
	Create(ctx context.Context, user *model.Users) error
	GetByEID(ctx context.Context, eid string) (*model.Users, error)
	GetByEIDUnscoped(ctx context.Context, eid string) (*model.Users, error)
}

type userService struct {
	store   store.Store
	storage storage.Storage
}

var _ UserSrv = &userService{}

func newUsers(srv *service) *userService {
	return &userService{store: srv.store, storage: srv.storage}
}

func (u userService) Create(ctx context.Context, user *model.Users) error {
	db := u.store.DB()

	if _, err := u.store.User().Get(ctx, db, user.Phone, options.WithQuery("phone = ?")); err == nil {
		return errors.Code(code.ErrPhoneAlreadyExist, "phone already exists")
	}

	if _, err := u.store.User().Get(ctx, db, user.EID, options.WithQuery("eid = ?")); err == nil {
		return errors.Code(code.ErrUsernameAlreadyExist, "eid already exists")
	}

	if err := u.store.User().Create(ctx, db, user); err != nil {
		if match, _ := regexp.MatchString("duplicate key value violates unique constraint .*", err.Error()); match {
			return errors.Code(code.ErrUserAlreadyExist, err.Error())
		}

		return errors.Code(code.ErrDatabase, err.Error())
	}

	return nil
}

func (u userService) GetByEID(ctx context.Context, eid string) (*model.Users, error) {
	user := &model.Users{}
	if err := u.storage.HGetAll(ctx, fmt.Sprintf(storage.KeyUser, eid), user); err == nil {
		log.Debugf("get user from storage: %+v", user.EID)
		return user, nil
	}

	db := u.store.DB()

	user, err := u.store.User().Get(ctx, db, eid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Code(code.ErrUserNotExist, err.Error())
		}
		return nil, errors.Code(code.ErrDatabase, err.Error())
	}

	_ = u.storage.HSetAllWithExpire(ctx, fmt.Sprintf(storage.KeyUser, eid), user, time.Hour)
	return user, nil
}

func (u userService) GetByEIDUnscoped(ctx context.Context, eid string) (*model.Users, error) {
	user := &model.Users{}
	if err := u.storage.HGetAll(ctx, fmt.Sprintf(storage.KeyUserUnscoped, eid), user); err == nil {
		log.Debugf("get user from storage: %+v", user.EID)
		return user, nil
	}

	db := u.store.DB()

	user, err := u.store.User().Get(ctx, db, eid, options.WithUnscoped())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Code(code.ErrUserNotExist, err.Error())
		}
		return nil, errors.Code(code.ErrDatabase, err.Error())
	}

	_ = u.storage.HSetAllWithExpire(ctx, fmt.Sprintf(storage.KeyUserUnscoped, eid), user, time.Hour)
	return user, nil
}
