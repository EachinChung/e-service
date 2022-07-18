package mysql

import (
	"fmt"
	"sync"

	"gorm.io/gorm"

	"github.com/eachinchung/component-base/options"
	"github.com/eachinchung/errors"

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

var (
	mysqlFactory store.Store
	once         sync.Once
)

// GetMySQLFactoryOr 使用给定的配置创建 mysql 工厂。
func GetMySQLFactoryOr(opts *options.MySQLOptions) (store.Store, error) {
	if opts == nil && mysqlFactory == nil {
		return nil, fmt.Errorf("获取 mysql 工厂失败")
	}

	var err error
	var dbIns *gorm.DB
	once.Do(func() {
		dbIns, err = opts.NewClient()
		mysqlFactory = &datastore{dbIns}
	})

	if mysqlFactory == nil || err != nil {
		return nil, fmt.Errorf("获取 mysql 工厂失败, mysqlFactory: %+v, error: %w", mysqlFactory, err)
	}

	return mysqlFactory, nil
}
