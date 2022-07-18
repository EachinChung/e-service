package user

import (
	"github.com/gin-gonic/gin"

	"github.com/eachinchung/component-base/core"
	"github.com/eachinchung/errors"
	"github.com/eachinchung/log"

	"github.com/eachinchung/e-service/internal/pkg/code"
	"github.com/eachinchung/e-service/internal/pkg/validator"
)

type GetUri struct {
	UserID uint64 `uri:"id" binding:"required"`
}

// Get get a user by the user identifier.
func (u *Controller) Get(c *gin.Context) {
	uri := &GetUri{}
	if err := c.ShouldBindUri(uri); err != nil {
		core.WriteResponse(
			c,
			validator.ParseValidationError(err),
			core.WithError(errors.Code(code.ErrValidation, err.Error())),
		)
		return
	}

	user, err := u.srv.Users().GetByUserID(c, uri.UserID)
	if err != nil {
		log.Errorf("get user error: %+v", err)
		core.WriteResponse(c, nil, core.WithError(err))
		return
	}

	core.WriteResponse(c, user)
}
