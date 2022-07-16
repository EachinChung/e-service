package code

// Common: basic errors.
const (
	// ErrSuccess - 200: success.
	ErrSuccess int = iota + 100001

	// ErrUnknown - 500: 服务器内部错误.
	ErrUnknown

	// ErrBind - 400: 请求参数格式错误.
	ErrBind

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
)

// common: 网络相关错误
const (
	// ErrNetworkRequest - 500: 网络请求错误.
	ErrNetworkRequest int = iota + 100201

	// ErrNetworkTimeOut - 500: 网络请求超时.
	ErrNetworkTimeOut
)
