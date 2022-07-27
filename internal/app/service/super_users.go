package service

import (
	"context"
	"fmt"
	"time"

	"github.com/eachinchung/e-service/internal/app/storage"
	"github.com/eachinchung/e-service/internal/app/store"
)

type SuperUsersSrv interface {
	Exists(ctx context.Context, eid string) (bool, error)
}

type superUserService struct {
	store   store.Store
	storage storage.Storage
}

var _ SuperUsersSrv = &superUserService{}

func newSuperUsers(srv *service) *superUserService {
	return &superUserService{store: srv.store, storage: srv.storage}
}

func (s superUserService) Exists(ctx context.Context, eid string) (bool, error) {
	if exists, err := s.storage.GetBool(ctx, fmt.Sprintf(storage.KeyIsSuperUser, eid)); err == nil {
		return exists, err
	}
	db := s.store.DB()

	exist, err := s.store.SuperUsers().Exist(ctx, db, eid)
	if err != nil {
		return exist, err
	}

	_ = s.storage.Set(ctx, fmt.Sprintf(storage.KeyIsSuperUser, eid), exist, time.Hour)
	return exist, nil
}
