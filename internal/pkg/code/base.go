package code

//go:generate codegen -type=int

// Common: basic errors.
const (
	// ErrUnknown - 500: 服务器内部错误.
	ErrUnknown int = iota + 100001

	// ErrValidation - 400: 参数验证失败.
	ErrValidation

	// ErrPageNotFound - 404: 资源不存在.
	ErrPageNotFound
)

// common: 授权和身份验证错误。
const (
	// ErrPermissionDenied - 403: 没有权限.
	ErrPermissionDenied int = iota + 100101

	// ErrNeedCaptcha - 403: 请验证您不是机器人.
	ErrNeedCaptcha

	// ErrCaptchaVerifyFailed - 400: 验证码校验失败.
	ErrCaptchaVerifyFailed

	// ErrCaptchaBusy - 500: 验证码服务繁忙, 请稍后再试.
	ErrCaptchaBusy

	// ErrFailedTokenCreation - 500: token 创建失败.
	ErrFailedTokenCreation

	// ErrEmptyToken - 401: 没有携带 token.
	ErrEmptyToken

	// ErrInvalidToken - 401: Token 无效.
	ErrInvalidToken

	// ErrExpiredToken - 401: token 已过期.
	ErrExpiredToken

	// ErrFailedAuthentication - 400: 账号或密码不正确.
	ErrFailedAuthentication

	// ErrNetworkUnsafe - 403: 网络环境不安全.
	ErrNetworkUnsafe
)

// common: database errors.
const (
	// ErrDatabase - 500: 数据库错误.
	ErrDatabase int = iota + 100201
)

// common: 用户相关错误
const (
	// ErrUserAlreadyExist - 400: 用户已存在.
	ErrUserAlreadyExist int = iota + 100301

	// ErrUserNotExist - 404: 用户不存在.
	ErrUserNotExist

	// ErrUsernameAlreadyExist - 400: 用户名已存在.
	ErrUsernameAlreadyExist

	// ErrPhoneAlreadyExist - 400: 该手机号码已注册.
	ErrPhoneAlreadyExist

	// ErrEmailAlreadyExist - 400: 该邮箱已注册.
	ErrEmailAlreadyExist

	// ErrUserStatusIsAbnormal - 403: 用户状态异常.
	ErrUserStatusIsAbnormal
)
