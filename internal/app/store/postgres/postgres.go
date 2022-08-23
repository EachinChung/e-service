package postgres

import (
	"fmt"
	"sync"

	"github.com/eachinchung/component-base/options"
	"github.com/eachinchung/errors"
	"gorm.io/gorm"

	"github.com/eachinchung/e-service/internal/app/store"
)

type datastore struct {
	db *gorm.DB
}

func (ds *datastore) DB() *gorm.DB {
	return ds.db
}

func (ds *datastore) Close() error {
	dbIns, err := ds.db.DB()
	if err != nil {
		return errors.Wrap(err, "获取 gorm db 实例失败")
	}

	return dbIns.Close()
}

func (ds *datastore) User() store.UserStore {
	return newUser()
}

func (ds *datastore) SuperUsers() store.SuperUsersStore {
	return newSuperUser()
}

var (
	factory store.Store
	once    sync.Once
)

// GetPostgresFactoryOr 使用给定的配置创建 postgres 工厂。
func GetPostgresFactoryOr(opts *options.PostgresOptions) (store.Store, error) {
	if opts == nil && factory == nil {
		return nil, fmt.Errorf("获取 postgres 工厂失败")
	}

	var err error
	var dbIns *gorm.DB
	once.Do(func() {
		dbIns, err = opts.NewClient()
		factory = &datastore{dbIns}
	})

	if factory == nil || err != nil {
		return nil, fmt.Errorf("获取 postgres 工厂失败, factory: %+v, error: %w", factory, err)
	}

	return factory, nil
}
