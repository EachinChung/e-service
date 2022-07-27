package casbin

import (
	"context"
	"fmt"
	"sync"

	"github.com/casbin/casbin/v2"
	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"

	"github.com/eachinchung/component-base/core"
	"github.com/eachinchung/component-base/middleware/auth"
	"github.com/eachinchung/component-base/options"
	"github.com/eachinchung/errors"
	"github.com/eachinchung/log"

	"github.com/eachinchung/e-service/internal/app/service"
	"github.com/eachinchung/e-service/internal/app/storage"
	"github.com/eachinchung/e-service/internal/app/store"
	"github.com/eachinchung/e-service/internal/app/store/model"
	"github.com/eachinchung/e-service/internal/pkg/code"
)

var (
	enforcer *casbin.Enforcer
	once     sync.Once
)

func GetEnforcerOr(opts *options.CasbinOptions) (*casbin.Enforcer, error) {
	var a *adapter.Adapter
	var err error

	once.Do(func() {
		s := store.Client()
		srv := service.NewService(s, storage.Client())

		if a, err = adapter.NewAdapterByDB(s.DB()); err != nil {
			return
		}
		if enforcer, err = casbin.NewEnforcer(opts.Model, a); err != nil {
			return
		}
		if err = enforcer.LoadPolicy(); err != nil {
			return
		}

		enforcer.AddFunction("isSuperUser", func(arguments ...any) (any, error) {
			rSub := arguments[0].(string)
			return srv.SuperUser().Exists(context.Background(), rSub)
		})
	})

	if err != nil {
		return nil, errors.Wrap(err, "获取 casbin enforcer 失败")
	}

	return enforcer, nil
}

func RBACMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := auth.ExtractClaimsFromContext(c)
		eid := claims["sub"].(string)
		srv := service.NewService(store.Client(), storage.Client())
		user, err := srv.Users().GetByEID(c, eid)
		if err != nil {
			log.L(c).Errorf("获取用户信息失败: %+v", err)
			core.WriteResponse(
				c,
				nil,
				core.WithError(errors.Code(code.ErrDatabase, "获取用户信息失败")),
				core.WithAbort(),
			)
			return
		}

		if user.State != model.StatusNormal {
			core.WriteResponse(
				c,
				nil,
				core.WithError(errors.Code(code.ErrUserStatusIsAbnormal, "用户已被禁用")),
				core.WithMessage(fmt.Sprintf("用户已被%s", user.State.Msg())),
				core.WithAbort(),
			)
			return
		}

		ok, err := Enforce(c, user.EID, c.Request.URL.Path, c.Request.Method)
		if err != nil {
			log.L(c).Errorf("获取用户权限失败: %+v", err)
			core.WriteResponse(
				c,
				nil,
				core.WithError(errors.Code(code.ErrDatabase, "获取用户权限失败")),
				core.WithAbort(),
			)
			return
		}

		if !ok {
			log.L(c).Warnf("用户 %s 没有权限: %+v", user.EID, c.Request.URL.Path)
			core.WriteResponse(
				c,
				nil,
				core.WithError(errors.Code(code.ErrPermissionDenied, "用户没有权限")),
				core.WithAbort(),
			)
			return
		}

		user.SaveToContext(c)
	}
}

//goland:noinspection SpellCheckingInspection
func Enforce(ctx context.Context, user any, permission ...any) (bool, error) {
	ok, err := enforcer.Enforce(joinSlice(user, permission...)...)
	if err != nil {
		return false, errors.Wrap(err, "获取用户权限失败")
	}

	log.L(ctx).Infof("用户 %s 校验权限: %+v 结果: %+v", user, permission, ok)
	return ok, nil
}

// AddPermissionForUser 添加用户权限
func AddPermissionForUser(ctx context.Context, user string, permission ...string) error {
	_, err := enforcer.AddPermissionForUser(user, permission...)
	if err != nil {
		return errors.Wrap(err, "添加用户权限失败")
	}
	log.L(ctx).Infof("用户 %s 添加权限: %+v", user, permission)
	return err
}

// DeletePermissionForUser 删除用户权限
func DeletePermissionForUser(ctx context.Context, user string, permission ...string) error {
	_, err := enforcer.DeletePermissionForUser(user, permission...)
	if err != nil {
		return errors.Wrap(err, "删除用户权限失败")
	}
	log.L(ctx).Infof("用户 %s 删除权限: %+v", user, permission)
	return err
}

// GetPermissionsForUser 获取用户或角色的权限
func GetPermissionsForUser(ctx context.Context, user string, domain ...string) [][]string {
	ps := enforcer.GetPermissionsForUser(user, domain...)
	log.L(ctx).Infof("用户 %s 权限: %+v", user, ps)
	return ps
}

// HasPermissionForUser 确定用户是否具有权限
func HasPermissionForUser(ctx context.Context, user string, permission ...string) bool {
	ok := enforcer.HasPermissionForUser(user, permission...)
	log.L(ctx).Infof("确定用户 %s 是否具有权限: %+v 结果: %+v", user, permission, ok)
	return ok
}

// joinSlice joins an any and a slice into a new slice.
func joinSlice(a any, b ...any) []any {
	res := make([]any, 0, len(b)+1)

	res = append(res, a)
	res = append(res, b...)
	return res
}
