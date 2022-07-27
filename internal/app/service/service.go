package service

import (
	"github.com/eachinchung/e-service/internal/app/storage"
	"github.com/eachinchung/e-service/internal/app/store"
)

// Service defines functions used to return resource interface.
type Service interface {
	Users() UserSrv
	SuperUser() SuperUsersSrv
}

type service struct {
	store   store.Store
	storage storage.Storage
}

// NewService returns Service interface.
func NewService(store store.Store, storage storage.Storage) Service {
	return &service{
		store:   store,
		storage: storage,
	}
}

func (s *service) Users() UserSrv {
	return newUsers(s)
}

func (s *service) SuperUser() SuperUsersSrv {
	return newSuperUsers(s)
}
