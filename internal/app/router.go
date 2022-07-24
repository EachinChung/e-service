package app

import (
	"github.com/eachinchung/e-service/internal/pkg/casbin"
	"github.com/gin-gonic/gin"

	"github.com/eachinchung/component-base/core"
	"github.com/eachinchung/errors"

	"github.com/eachinchung/e-service/internal/app/controller/v1/user"
	"github.com/eachinchung/e-service/internal/app/store/postgres"
	"github.com/eachinchung/e-service/internal/pkg/code"
)

func initRouter(g *gin.Engine) {
	installController(g)
}

func installController(g *gin.Engine) {
	jwtStrategy := newJWTAuth()

	auth := g.Group("/auth")
	{
		auth.POST("token", jwtStrategy.LoginHandler)
		auth.PUT("token", jwtStrategy.RefreshHandler)
	}

	g.NoRoute(jwtStrategy.MiddlewareFunc(), func(c *gin.Context) {
		core.WriteResponse(c, nil, core.WithError(errors.Code(code.ErrPageNotFound, "page not found")))
	})

	storeIns, _ := postgres.GetPostgresFactoryOr(nil)
	v1 := g.Group("/v1")
	{
		users := v1.Group("/users")
		{
			userController := user.NewController(storeIns)

			users.Use(jwtStrategy.MiddlewareFunc(), casbin.RBACMiddleWare())
			users.POST("", userController.Create)
			users.GET(":username", userController.Get)
		}
	}
}
