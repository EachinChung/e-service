package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var ErrKeyNotFound = fmt.Errorf("key not found")

type Storage interface {
	RDB() *redis.Client
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	GetBool(ctx context.Context, key string) (bool, error)

	HSet(ctx context.Context, key string, values ...any) error
	HSetAllWithExpire(ctx context.Context, key string, model any, expiration time.Duration) error
	HMGet(ctx context.Context, key string, fields ...string) ([]any, error)
	HGetAll(ctx context.Context, key string, model any) error

	Expire(ctx context.Context, key string, expiration time.Duration) error
}

// Client 返回 store 客户端实例。
func Client() Storage {
	if rs == nil {
		panic("store client is not set")
	}
	return rs
}
