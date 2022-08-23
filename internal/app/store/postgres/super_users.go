package postgres

import (
	"context"

	"github.com/eachinchung/errors"
	"gorm.io/gorm"

	"github.com/eachinchung/e-service/internal/app/store"
	"github.com/eachinchung/e-service/internal/app/store/model"
)

type superUser struct{}

func newSuperUser() *superUser {
	return &superUser{}
}

var _ store.SuperUsersStore = &superUser{}

func (s superUser) Exist(ctx context.Context, db *gorm.DB, eid string) (bool, error) {
	var count int64
	if err := db.Model(&model.SuperUsers{}).Where("eid = ?", eid).Count(&count).Error; err != nil {
		return false, errors.Wrap(err, "failed to check super user")
	}
	return count > 0, nil
}
