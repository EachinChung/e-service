package options

import (
	"fmt"

	"github.com/spf13/pflag"

	"github.com/eachinchung/e-service/internal/pkg/server"
)

// ServerRunOptions 运行服务器时的选项
type ServerRunOptions struct {
	Mode           string   `json:"mode"            mapstructure:"mode"`
	Healthz        bool     `json:"healthz"         mapstructure:"healthz"`
	Middlewares    []string `json:"middlewares"     mapstructure:"middlewares"`
	TrustedProxies []string `json:"trusted-proxies" mapstructure:"trusted-proxies"`
	BindAddress    string   `json:"bind-address"    mapstructure:"bind-address"`
	BindPort       int      `json:"bind-port"       mapstructure:"bind-port"`
}

// NewServerRunOptions 创建一个带有默认参数的 ServerRunOptions 对象。
func NewServerRunOptions() *ServerRunOptions {
	defaults := server.NewConfig()

	return &ServerRunOptions{
		Mode:           defaults.Mode,
		Healthz:        defaults.Healthz,
		Middlewares:    defaults.Middlewares,
		TrustedProxies: defaults.TrustedProxies,
		BindAddress:    "127.0.0.1",
		BindPort:       8080,
	}
}

// ApplyTo 应用配置
func (s *ServerRunOptions) ApplyTo(c *server.Config) error {
	c.Mode = s.Mode
	c.Healthz = s.Healthz
	c.Middlewares = s.Middlewares
	c.TrustedProxies = s.TrustedProxies

	c.Serving = &server.ServingInfo{
		BindAddress: s.BindAddress,
		BindPort:    s.BindPort,
	}

	return nil
}

// Validate 验证选项字段。
func (s *ServerRunOptions) Validate() []error {
	var errors []error

	if s.BindPort < 1 || s.BindPort > 65535 {
		errors = append(errors,
			fmt.Errorf("--insecure.bind-port %v 必须介于 1 和 65535 之间，包括 1 和 65535。", s.BindPort))
	}

	return errors
}

// AddFlags 将 ServerRun 的各个字段追加到传入的 pflag.FlagSet 变量中。
func (s *ServerRunOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(
		&s.Mode,
		"server.mode",
		s.Mode,
		"以指定的服务器模式启动服务器。支持的服务器模式: debug, test, release。",
	)

	fs.BoolVar(
		&s.Healthz,
		"server.healthz",
		s.Healthz,
		"是否开启健康检查，如果开启会安装 /healthz 路由，默认 true",
	)

	fs.StringSliceVar(
		&s.Middlewares,
		"server.middlewares",
		s.Middlewares,
		"加载的 gin 中间件列表，多个中间件，逗号(,)隔开",
	)

	fs.StringSliceVar(
		&s.TrustedProxies,
		"server.trusted-proxies",
		s.Middlewares,
		"受信任的代理地址，多个地址，逗号(,)隔开",
	)

	fs.StringVar(
		&s.BindAddress,
		"server.bind-address",
		s.BindAddress,
		"绑定 IP 地址，设置为 0.0.0.0 表示使用全部网络接口，默认为 127.0.0.1",
	)

	fs.IntVar(&s.BindPort, "server.bind-port", s.BindPort, "监听端口，默认为 8080")
}
