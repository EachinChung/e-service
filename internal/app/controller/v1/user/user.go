package user

import (
	"github.com/eachinchung/e-service/internal/app/service"
	"github.com/eachinchung/e-service/internal/app/store"
)

// Controller create a user handler used to handle request for user resource.
type Controller struct {
	srv service.Service
}

// NewController creates a user handler.
func NewController(store store.Store) *Controller {
	return &Controller{
		srv: service.NewService(store),
	}
}
