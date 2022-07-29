package user

import (
	"github.com/eachinchung/e-service/internal/app/validator"
	"github.com/gin-gonic/gin"

	"github.com/eachinchung/component-base/auth"
	"github.com/eachinchung/component-base/core"
	"github.com/eachinchung/component-base/utils/idutil"
	"github.com/eachinchung/errors"
	"github.com/eachinchung/log"

	"github.com/eachinchung/e-service/internal/app/store/model"
	"github.com/eachinchung/e-service/internal/pkg/code"
)

type createBody struct {
	Phone        string  `json:"phone"         binding:"required,len=11,phone"`                  // 手机号
	Nickname     string  `json:"nickname"      binding:"required,min=1,max=32"`                  // 昵称
	EID          *string `json:"eid"           binding:"omitempty,min=6,max=20,eid,is_not_role"` // 用户名
	Password     string  `json:"password"      binding:"required,min=6,password"`                // 密码
	ActivateCode string  `json:"activate_code" binding:"required,min=6"`                         // 激活码
	Captcha      string  `json:"captcha"       binding:"required,len=4"`                         // 验证码
}

func (u *Controller) Create(c *gin.Context) {
	body := &createBody{}
	if err := c.ShouldBindJSON(body); err != nil {
		core.WriteResponse(
			c,
			validator.ParseValidationError(err),
			core.WithError(errors.Code(code.ErrValidation, err.Error())),
		)
		return
	}

	pwdHash, _ := auth.HashPassword(body.Password)

	user := &model.Users{
		Phone:        body.Phone,
		Nickname:     body.Nickname,
		PasswordHash: pwdHash,
	}

	if body.EID == nil {
		user.EID = idutil.GetInstanceID(idutil.GenUint64ID(), "eid")
	} else {
		user.EID = *body.EID
	}

	if err := u.srv.Users().Create(c, user); err != nil {
		if !errors.IsCode(err, code.ErrEmailAlreadyExist) &&
			!errors.IsCode(err, code.ErrPhoneAlreadyExist) &&
			!errors.IsCode(err, code.ErrUsernameAlreadyExist) {
			log.Errorf("create user error: %+v", err)
		}

		core.WriteResponse(c, nil, core.WithError(err))
		return
	}

	core.WriteResponse(c, user)
}
