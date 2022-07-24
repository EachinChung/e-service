package postgres

import (
	"context"

	"github.com/eachinchung/e-service/internal/app/store/model"
	"github.com/eachinchung/errors"

	"github.com/eachinchung/e-service/internal/app/store"
	"gorm.io/gorm"
)

type superUser struct{}

func newSuperUser() *superUser {
	return &superUser{}
}

var _ store.SuperUsersStore = &superUser{}

func (s superUser) Exist(ctx context.Context, db *gorm.DB, username string) (bool, error) {
	var count int64
	if err := db.Model(&model.SuperUsers{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, errors.Wrap(err, "failed to check super user")
	}
	return count > 0, nil
}
