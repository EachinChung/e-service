package app

import (
	"github.com/gin-gonic/gin"

	"github.com/eachinchung/component-base/core"
	"github.com/eachinchung/errors"

	"github.com/eachinchung/e-service/internal/pkg/code"
)

func initRouter(g *gin.Engine) {
	installController(g)
}

func installController(g *gin.Engine) {
	g.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, errors.Code(code.ErrPageNotFound, "Page not found."), nil)
	})
}
