package storage

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/eachinchung/component-base/options"
	"github.com/eachinchung/errors"
	"github.com/eachinchung/log"
	"github.com/go-redis/redis/v8"
)

var (
	once sync.Once
	rs   Storage
)

func GetRedisClientOr(opts *options.RedisOptions) (Storage, error) {
	var err error

	once.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", opts.Host, opts.Port),
			Password: opts.Password,
			DB:       opts.Database,
		})

		_, err = client.Ping(context.Background()).Result()
		rs = newRedisStorage(client)
	})

	if err != nil {
		return nil, err
	}

	return rs, nil
}

type redisStorage struct {
	client *redis.Client
}

func newRedisStorage(client *redis.Client) *redisStorage {
	return &redisStorage{client}
}

func (r *redisStorage) RDB() *redis.Client {
	return r.client
}

func (r *redisStorage) Get(ctx context.Context, key string) (string, error) {
	log.L(ctx).Debugf("[STORE] GET key is: %s", key)
	val, err := r.client.Get(ctx, key).Result()
	switch {
	case err == redis.Nil:
		return val, ErrKeyNotFound
	case err != nil:
		log.L(ctx).Errorf("[STORE] GET key is: %s, err: %+v", key, err)
		return val, err
	}
	return val, err
}

func (r *redisStorage) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	log.L(ctx).Debugf("[STORE] SET key is: %s", key)
	err := r.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		log.L(ctx).Errorf("[STORE] SET key is: %s, err: %+v", key, err)
		return err
	}
	return err
}

func (r *redisStorage) GetBool(ctx context.Context, key string) (bool, error) {
	log.L(ctx).Debugf("[STORE] GET bool key is: %s", key)
	val, err := r.client.Get(ctx, key).Bool()
	switch {
	case err == redis.Nil:
		return val, ErrKeyNotFound
	case err != nil:
		log.L(ctx).Errorf("[STORE] GET bool key is: %s, err: %+v", key, err)
		return val, err
	}
	return val, err
}

func (r *redisStorage) GetInt(ctx context.Context, key string) (int, error) {
	log.L(ctx).Debugf("[STORE] GET int key is: %s", key)

	val, err := r.client.Get(ctx, key).Int()
	switch {
	case err == redis.Nil:
		return val, ErrKeyNotFound
	case err != nil:
		log.L(ctx).Errorf("[STORE] GET int key is: %s, err: %+v", key, err)
		return val, err
	}
	return val, err
}

func (r *redisStorage) GetInt64(ctx context.Context, key string) (int64, error) {
	log.L(ctx).Debugf("[STORE] GET int64 key is: %s", key)

	val, err := r.client.Get(ctx, key).Int64()
	switch {
	case err == redis.Nil:
		return val, ErrKeyNotFound
	case err != nil:
		log.L(ctx).Errorf("[STORE] GET int64 key is: %s, err: %+v", key, err)
		return val, err
	}
	return val, err
}

func (r *redisStorage) GetUint64(ctx context.Context, key string) (uint64, error) {
	log.L(ctx).Debugf("[STORE] GET uint64 key is: %s", key)

	val, err := r.client.Get(ctx, key).Uint64()
	switch {
	case err == redis.Nil:
		return val, ErrKeyNotFound
	case err != nil:
		log.L(ctx).Errorf("[STORE] GET uint64 key is: %s, err: %+v", key, err)
		return val, err
	}
	return val, err
}

func (r *redisStorage) Incr(ctx context.Context, key string) (int64, error) {
	log.L(ctx).Debugf("[STORE] INCR key is: %s", key)
	val, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		log.L(ctx).Errorf("[STORE] INCR key is: %s, err: %+v", key, err)
		return val, err
	}
	return val, nil
}

func (r *redisStorage) HSet(ctx context.Context, key string, values ...any) error {
	log.L(ctx).Debugf("[STORE] HSET key is: %s", key)
	err := r.client.HSet(ctx, key, values...).Err()
	if err != nil {
		log.L(ctx).Errorf("[STORE] HSET key is: %s, err: %+v", key, err)
		return err
	}
	return nil
}

func (r *redisStorage) HSetAllWithExpire(
	ctx context.Context,
	key string,
	model any,
	expiration time.Duration,
) error {
	log.L(ctx).Debugf("[STORE] HSETALL key is: %s, Expire: %v", key, expiration)

	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Ptr {
		return errors.Errorf("model only accepts struct pointer: got %T", v)
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return errors.Errorf("model only accepts struct pointer: got %T", v)
	}

	t := v.Type()

	if _, err := r.client.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		for i := 0; i < v.NumField(); i++ {
			var k string
			tagValue := t.Field(i).Tag.Get("redis")
			switch tagValue {
			case "-":
				continue
			case "":
				k = t.Field(i).Name
			default:
				k = tagValue
			}

			var val any
			switch t.Field(i).Type.Kind() {
			case reflect.String:
				val = v.Field(i).String()
			case reflect.Int8, reflect.Int16, reflect.Int, reflect.Int32, reflect.Int64:
				val = v.Field(i).Int()
			case reflect.Uint8, reflect.Uint, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				val = v.Field(i).Uint()
			case reflect.Float32, reflect.Float64:
				val = v.Field(i).Float()
			case reflect.Bool:
				val = v.Field(i).Bool()
			case reflect.Struct:
				if scan := v.Field(i).Addr().MethodByName("Value"); scan.IsValid() {
					val = scan.Call([]reflect.Value{})[0].Interface()
				} else {
					val = v.Field(i).Interface()
				}
			default:
				val = v.Field(i).Interface()
			}

			rdb.HSet(ctx, key, k, val)
		}

		if expiration > 0 {
			rdb.Expire(ctx, key, expiration)
		}
		return nil
	}); err != nil {
		log.L(ctx).Errorf("[STORE] HSETALL key is: %s, err: %+v", key, err)
		return err
	}
	return nil
}

func (r *redisStorage) HMGet(ctx context.Context, key string, fields ...string) ([]any, error) {
	log.L(ctx).Debugf("[STORE] HMGET key is: %s", key)
	result, err := r.client.HMGet(ctx, key, fields...).Result()
	switch {
	case err == redis.Nil:
		return nil, ErrKeyNotFound
	case err != nil:
		log.L(ctx).Errorf("[STORE] HMGET key is: %s, err: %+v", key, err)
		return nil, err
	}
	return result, nil
}

func (r *redisStorage) HGetAll(ctx context.Context, key string, model any) error {
	log.L(ctx).Debugf("[STORE] HGETALL key is: %s", key)

	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Ptr {
		return errors.Errorf("model only accepts struct pointer: got %T", v)
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return errors.Errorf("model only accepts struct pointer: got %T", v)
	}

	result, err := r.client.HGetAll(ctx, key).Result()
	switch {
	case err == redis.Nil:
		return ErrKeyNotFound
	case err != nil:
		log.L(ctx).Errorf("[STORE] HGETALL key is: %s, err: %+v", key, err)
		return err
	case len(result) == 0:
		return ErrKeyNotFound
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		var key string
		tagValue := t.Field(i).Tag.Get("redis")
		switch tagValue {
		case "-":
			continue
		case "":
			key = t.Field(i).Name
		default:
			key = tagValue
		}

		strVal, ok := result[key]
		if !ok {
			continue
		}

		if strVal == "" {
			continue
		}

		switch t.Field(i).Type.Kind() {
		case reflect.String:
			v.Field(i).SetString(strVal)
		case reflect.Int8, reflect.Int16, reflect.Int, reflect.Int32, reflect.Int64:
			val, err := strconv.ParseInt(strVal, 10, 64)
			if err != nil {
				return err
			}
			v.Field(i).SetInt(val)
		case reflect.Uint8, reflect.Uint, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			val, err := strconv.ParseUint(strVal, 10, 64)
			if err != nil {
				return err
			}
			v.Field(i).SetUint(val)
		case reflect.Float32, reflect.Float64:
			val, err := strconv.ParseFloat(strVal, 64)
			if err != nil {
				return err
			}
			v.Field(i).SetFloat(val)
		case reflect.Bool:
			val, err := strconv.ParseBool(strVal)
			if err != nil {
				return err
			}
			v.Field(i).SetBool(val)
		case reflect.Struct:
			switch t.Field(i).Type.String() {
			case "time.Time":
				tmp, err := time.ParseInLocation(time.RFC3339, strVal, time.Local)
				if err != nil {
					return err
				}
				v.Field(i).Set(reflect.ValueOf(tmp))
			default:
				scan := v.Field(i).Addr().MethodByName("Scan")

				if !scan.IsValid() {
					err = errors.Errorf("unknown type, name: %+v, kind: %s", t.Field(i).Name, t.Field(i).Type.Kind())
					log.L(ctx).Errorf("[STORE] HGETALL unsupported Scan, key is: %s, err: %+v", key, err)
					return err
				}

				scan.Call([]reflect.Value{reflect.ValueOf(strVal)})
			}
		default:
			err = errors.Errorf("unknown type, name: %+v, kind: %s", t.Field(i).Name, t.Field(i).Type.Kind())
			log.L(ctx).Errorf("[STORE] HGETALL key is: %s, err: %+v", key, err)
			return err
		}
	}

	return nil
}

func (r *redisStorage) Expire(ctx context.Context, key string, expiration time.Duration) error {
	log.L(ctx).Debugf("[STORE] EXPIRE key is: %s", key)
	err := r.client.Expire(ctx, key, expiration).Err()
	if err != nil {
		log.L(ctx).Errorf("[STORE] EXPIRE key is: %s, err: %+v", key, err)
		return err
	}
	return nil
}
