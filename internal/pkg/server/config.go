package server

import (
	"net"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Config 是用于配置 APIServer 的结构体。
type Config struct {
	Serving        *ServingInfo
	Mode           string
	Middlewares    []string
	TrustedProxies []string
	Healthz        bool
}

// NewConfig 返回一个具有默认值的 Config 结构体。
func NewConfig() *Config {
	return &Config{
		Healthz: true,
		Mode:    gin.ReleaseMode,
	}
}

// ServingInfo server 的运行配置.
type ServingInfo struct {
	BindAddress string
	BindPort    int
}

// Address 将主机IP地址和主机端口号加入到一个地址字符串中，如: 127.0.0.1:8080。
func (s *ServingInfo) Address() string {
	return net.JoinHostPort(s.BindAddress, strconv.Itoa(s.BindPort))
}

// CompletedConfig 是 APIServer 的完整配置.
type CompletedConfig struct {
	*Config
}

// Complete 填写任何未设置的字段，这些字段需要具有有效数据并且可以从其他字段派生。
// 如果您要使用“ApplyOptions”，请先执行此操作。 它正在改变接收器。
func (c *Config) Complete() CompletedConfig {
	return CompletedConfig{c}
}

// New 从给定的配置中返回一个新的 APIServer 实例。
func (c CompletedConfig) New() (*APIServer, error) {
	s := &APIServer{
		ServingInfo:    c.Serving,
		mode:           c.Mode,
		healthz:        c.Healthz,
		middlewares:    c.Middlewares,
		trustedProxies: c.TrustedProxies,
		Engine:         gin.New(),
	}

	initGenericAPIServer(s)

	return s, nil
}
