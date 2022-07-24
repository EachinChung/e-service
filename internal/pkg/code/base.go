package code

//go:generate codegen -type=int

// Common: basic errors.
const (
	// ErrSuccess - 200: success.
	ErrSuccess int = iota + 100001

	// ErrUnknown - 500: 服务器内部错误.
	ErrUnknown

	// ErrValidation - 400: 参数验证失败.
	ErrValidation

	// ErrPageNotFound - 404: 资源不存在.
	ErrPageNotFound
)

// common: 授权和身份验证错误。
const (
	// ErrTokenInvalid - 401: Token 不合法.
	ErrTokenInvalid int = iota + 100101

	// ErrInvalidAuthHeader - 401: Authorization 不合法.
	ErrInvalidAuthHeader

	// ErrMissingAuthHeader - 401: Authorization 是空的.
	ErrMissingAuthHeader

	// ErrUsernameOrPasswordIncorrect - 401: 账号或密码错误.
	ErrUsernameOrPasswordIncorrect

	// ErrPermissionDenied - 403: 没有权限.
	ErrPermissionDenied

	// ErrFailedTokenCreation - 401: token 创建失败.
	ErrFailedTokenCreation

	// ErrExpiredToken - 401: token 已过期, 无法刷新.
	ErrExpiredToken

	// ErrMissingExpField - 400: 缺少 exp 字段.
	ErrMissingExpField

	// ErrWrongFormatOfExp - 400: exp 必须是 float64 格式.
	ErrWrongFormatOfExp

	// ErrEmptyToken - 401: 没有携带 token.
	ErrEmptyToken

	// ErrInvalidSigningAlgorithm - 400: 无效签名算法.
	ErrInvalidSigningAlgorithm

	// ErrFailedAuthentication - 401: 用户名或密码不正确.
	ErrFailedAuthentication
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
