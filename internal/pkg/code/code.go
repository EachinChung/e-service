package code

import (
	"net/http"

	"github.com/eachinchung/errors"
	"github.com/novalagung/gubrak/v2"
)

type ErrCode struct {
	// C 指的是 ErrCode 的错误码。
	C int

	// HTTP 关联 ErrCode 的 HTTP 状态码。
	HTTP int

	// Ext 外部 (用户) 面对的错误文本。
	Ext string
}

var _ errors.Coder = &ErrCode{}

// Code 返回 ErrCode 的错误代码。
func (coder ErrCode) Code() int {
	return coder.C
}

// String 返回外部 (用户) 面对的错误文本。
func (coder ErrCode) String() string {
	return coder.Ext
}

// HTTPStatus 返回关联的 HTTP 状态码。
func (coder ErrCode) HTTPStatus() int {
	if coder.HTTP == 0 {
		return http.StatusInternalServerError
	}

	return coder.HTTP
}

func register(code int, httpStatus int, message string) {
	found := gubrak.From([]int{200, 400, 401, 403, 404, 500}).Contains(httpStatus).Result()
	if !found {
		panic("http code not in `200, 400, 401, 403, 404, 500`")
	}

	coder := &ErrCode{
		C:    code,
		HTTP: httpStatus,
		Ext:  message,
	}

	errors.MustRegister(coder)
}
