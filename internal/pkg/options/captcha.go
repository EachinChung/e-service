package options

import (
	"github.com/spf13/pflag"
)

// CaptchaOptions jwt 配置选项
type CaptchaOptions struct {
	AppID        uint64 `json:"app-id"          mapstructure:"app-id"`
	AppSecretKey string `json:"app-secret-key"  mapstructure:"app-secret-key"`
	SecretID     string `json:"secret-id"       mapstructure:"secret-id"`
	SecretKey    string `json:"secret-key"      mapstructure:"secret-key"`
	IsVerify     bool   `json:"is-verify"       mapstructure:"is-verify"`
}

// NewCaptchaOptions 创建一个带有默认参数的 CaptchaOptions 对象。
func NewCaptchaOptions() *CaptchaOptions {
	return &CaptchaOptions{
		IsVerify: false,
	}
}

// Validate 验证选项字段。
func (s *CaptchaOptions) Validate() []error {
	return []error{}
}

// AddFlags 将 tencent-cloud 的各个字段追加到传入的 pflag.FlagSet 变量中。
func (s *CaptchaOptions) AddFlags(fs *pflag.FlagSet) {
	if fs == nil {
		return
	}

	fs.Uint64Var(
		&s.AppID,
		"captcha.app-id",
		s.AppID,
		"验证码应用ID",
	)
	fs.StringVar(
		&s.AppSecretKey,
		"captcha.app-secret-key",
		s.AppSecretKey,
		"验证码应用密钥",
	)
	fs.StringVar(
		&s.SecretID,
		"captcha.secret-id",
		s.SecretID,
		"验证码密钥ID",
	)
	fs.StringVar(
		&s.SecretKey,
		"captcha.secret-key",
		s.SecretKey,
		"验证码密钥",
	)
	fs.BoolVar(
		&s.IsVerify,
		"captcha.is-verify",
		s.IsVerify,
		"是否验证",
	)
}
