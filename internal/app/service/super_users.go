package service

import (
	"context"

	"github.com/eachinchung/e-service/internal/app/store"
)

type SuperUsersSrv interface {
	Exist(ctx context.Context, username string) (bool, error)
}

type superUserService struct {
	store store.Store
}

var _ SuperUsersSrv = &superUserService{}

func newSuperUsers(srv *service) *superUserService {
	return &superUserService{store: srv.store}
}

func (s superUserService) Exist(ctx context.Context, username string) (bool, error) {
	db := s.store.DB()
	return s.store.SuperUsers().Exist(ctx, db, username)
}
