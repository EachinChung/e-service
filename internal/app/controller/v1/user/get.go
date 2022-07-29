package user

import (
	"github.com/eachinchung/e-service/internal/app/store/casbin"
	"github.com/eachinchung/e-service/internal/app/validator"
	"github.com/gin-gonic/gin"

	"github.com/eachinchung/component-base/core"
	"github.com/eachinchung/errors"
	"github.com/eachinchung/log"

	"github.com/eachinchung/e-service/internal/app/store/model"
	"github.com/eachinchung/e-service/internal/pkg/code"
)

type getUri struct {
	EID string `uri:"eid" binding:"required"`
}

// GetByEID get a user by the user identifier.
func (u *Controller) GetByEID(c *gin.Context) {
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
	ok, err := casbin.Enforce(c, user.EID, "admin:user", "get")
	if err != nil {
		log.Errorf("get user error: %+v", err)
		core.WriteResponse(c, nil, core.WithError(errors.Code(code.ErrDatabase, err.Error())))
		return
	}

	if !ok {
		if user.EID != uri.EID {
			core.WriteResponse(c, nil, core.WithError(errors.Code(code.ErrPermissionDenied, "无权获取此用户")))
			return
		}

		core.WriteResponse(c, user)
		return
	}

	user, err = u.srv.Users().GetByEIDUnscoped(c, uri.EID)
	if err != nil {
		log.Errorf("get user error: %+v", err)
		core.WriteResponse(c, nil, core.WithError(err))
		return
	}
	core.WriteResponse(c, user.AdminResponse())
}
