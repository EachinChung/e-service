package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/eachinchung/component-base/auth"
	"github.com/eachinchung/component-base/core"
	"github.com/eachinchung/errors"
	"github.com/eachinchung/log"

	"github.com/eachinchung/e-service/internal/app/store/model"
	"github.com/eachinchung/e-service/internal/pkg/code"
	"github.com/eachinchung/e-service/internal/pkg/validator"
)

type CreateBody struct {
	Phone    string `json:"phone" binding:"required,len=11,phone"`             // 手机号
	Username string `json:"username" binding:"required,min=6,max=20,username"` // 用户名
	Password string `json:"password" binding:"required,min=6,password"`        // 密码
}

func (u *Controller) Create(c *gin.Context) {
	body := &CreateBody{}
	if err := c.ShouldBindJSON(body); err != nil {
		core.WriteResponse(
			c,
			validator.ParseValidationError(err),
			core.WithError(errors.Code(code.ErrValidation, err.Error())),
		)
		return
	}

	pwdHash, _ := auth.HashPassword(body.Password)

	if err := u.srv.Users().Create(c, &model.Users{
		Phone:        body.Phone,
		Username:     body.Username,
		PasswordHash: pwdHash,
	}); err != nil {
		log.Errorf("create user error: %+v", err)
		core.WriteResponse(c, nil, core.WithError(err))
		return
	}

	c.Status(http.StatusCreated)
}
