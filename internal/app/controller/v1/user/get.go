package user

import (
	"github.com/eachinchung/e-service/internal/app/store/model"
	"github.com/eachinchung/e-service/internal/pkg/casbin"
	"github.com/eachinchung/log"
	"github.com/gin-gonic/gin"

	"github.com/eachinchung/component-base/core"
	"github.com/eachinchung/e-service/internal/pkg/code"
	"github.com/eachinchung/e-service/internal/pkg/validator"
	"github.com/eachinchung/errors"
)

type getUri struct {
	Username string `uri:"username" binding:"required"`
}

// Get get a user by the user identifier.
func (u *Controller) Get(c *gin.Context) {
	uri := &getUri{}
	if err := c.ShouldBindUri(uri); err != nil {
		core.WriteResponse(
			c,
			validator.ParseValidationError(err),
			core.WithError(errors.Code(code.ErrValidation, err.Error())),
		)
		return
	}

	user := model.ExtractUsersFromContext(c)
	ok, err := casbin.Enforce(c, user.Username, "admin:user", "get")
	if err != nil {
		log.Errorf("get user error: %+v", err)
		core.WriteResponse(c, nil, core.WithError(errors.Code(code.ErrDatabase, err.Error())))
		return
	}

	if !ok {
		if user.Username != uri.Username {
			core.WriteResponse(c, nil, core.WithError(errors.Code(code.ErrPermissionDenied, "无权获取此用户")))
			return
		}

		core.WriteResponse(c, user)
		return
	}

	user, err = u.srv.Users().GetByUsername(c, uri.Username)
	if err != nil {
		log.Errorf("get user error: %+v", err)
		core.WriteResponse(c, nil, core.WithError(err))
		return
	}
	core.WriteResponse(c, user.Map())
}