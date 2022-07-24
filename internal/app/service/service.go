package service

import "github.com/eachinchung/e-service/internal/app/store"

// Service defines functions used to return resource interface.
type Service interface {
	Users() UserSrv
	SuperUser() SuperUsersSrv
}

type service struct {
	store store.Store
}

// NewService returns Service interface.
func NewService(store store.Store) Service {
	return &service{
		store: store,
	}
}

func (s *service) Users() UserSrv {
	return newUsers(s)
}

func (s *service) SuperUser() SuperUsersSrv {
	return newSuperUsers(s)
}
