package code

// init register error codes defines in this source code to `github.com/EachinChung/errors`
func init() {
	register(ErrSuccess, 200, "success")
	register(ErrUnknown, 500, "服务器内部错误")
	register(ErrBind, 400, "请求参数格式错误")
	register(ErrValidation, 400, "参数验证失败")
	register(ErrPageNotFound, 404, "资源不存在")
	register(ErrTokenInvalid, 401, "Token 不合法")
	register(ErrInvalidAuthHeader, 401, "Authorization 不合法")
	register(ErrMissingAuthHeader, 401, "Authorization 是空的")
	register(ErrUsernameOrPasswordIncorrect, 401, "账号或密码错误")
	register(ErrPermissionDenied, 403, "没有权限")
	register(ErrNetworkRequest, 500, "网络请求错误")
	register(ErrNetworkTimeOut, 500, "网络请求超时")
}
