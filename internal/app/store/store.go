package store

import (
	"sync"

	"gorm.io/gorm"
)

var (
	once   sync.Once
	client Store
)

type Store interface {
	DB() *gorm.DB
	Close() error

	User() UserStore
}

// Client 返回 store 客户端实例。
func Client() Store {
	if client == nil {
		panic("store client is not set")
	}
	return client
}

// SetClient 设置 store 客户端。
func SetClient(s Store) {
	once.Do(func() { client = s })
}
