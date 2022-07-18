package config

import (
	"sync"

	"github.com/eachinchung/e-service/internal/app/options"
)

type Config struct {
	*options.Options
}

var (
	cfg  *Config
	once sync.Once
)

// GetConfigIns 基于给定的命令行或配置文件选项创建一个配置实例。
func GetConfigIns(opts *options.Options) *Config {
	once.Do(func() { cfg = &Config{opts} })
	return cfg
}
