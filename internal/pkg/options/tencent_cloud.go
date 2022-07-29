package options

import (
	"github.com/spf13/pflag"
)

// TencentCloudOptions jwt 配置选项
type TencentCloudOptions struct {
	CaptchaAppID        uint64 `json:"captcha-app-id"          mapstructure:"captcha-app-id"`
	CaptchaAppSecretKey string `json:"captcha-app-secret-key"  mapstructure:"captcha-app-secret-key"`
}

// NewTencentCloudOptions 创建一个带有默认参数的 TencentCloudOptions 对象。
func NewTencentCloudOptions() *TencentCloudOptions {
	return &TencentCloudOptions{}
}

// Validate 验证选项字段。
func (s *TencentCloudOptions) Validate() []error {
	return []error{}
}

// AddFlags 将 tencent-cloud 的各个字段追加到传入的 pflag.FlagSet 变量中。
func (s *TencentCloudOptions) AddFlags(fs *pflag.FlagSet) {
	if fs == nil {
		return
	}

	fs.Uint64Var(
		&s.CaptchaAppID,
		"tencent-cloud.captcha-app-id",
		s.CaptchaAppID,
		"验证码应用ID",
	)
	fs.StringVar(
		&s.CaptchaAppSecretKey,
		"tencent-cloud.captcha-app-secret-key",
		s.CaptchaAppSecretKey,
		"验证码应用密钥",
	)
}
