package app

import (
	"strconv"
	"time"

	"github.com/eachinchung/component-base/core"
	"github.com/eachinchung/component-base/db/options"
	"github.com/eachinchung/component-base/middleware/auth"
	"github.com/eachinchung/component-base/utils/idutil"
	"github.com/eachinchung/component-base/verification"
	"github.com/eachinchung/errors"
	"github.com/eachinchung/log"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/eachinchung/e-service/internal/app/config"
	"github.com/eachinchung/e-service/internal/app/middleware"
	"github.com/eachinchung/e-service/internal/app/store"
	"github.com/eachinchung/e-service/internal/app/store/model"
	"github.com/eachinchung/e-service/internal/app/validator"
	"github.com/eachinchung/e-service/internal/pkg/code"
)

const (
	// APIServerAudience defines the value of jwt audience field.
	APIServerAudience = "api.service.eachin-life.com"

	// APIServerIssuer defines the value of jwt issuer field.
	APIServerIssuer = "e-service"
)

type loginInfo struct {
	Username string `form:"username" json:"username" binding:"required,max=20"`
	Password string `form:"password" json:"password" binding:"required"`
}

type tokenRsp struct {
	Token  string    `json:"token"`
	Expire time.Time `json:"expire"`
}

func newJWTAuth() *auth.GinJWTMiddleware {
	cfg := config.GetConfigIns(nil)

	jwtMiddleware, _ := auth.New(&auth.GinJWTMiddleware{
		Realm:            cfg.JWTOptions.Realm,
		SigningAlgorithm: "HS512",
		Key:              []byte(cfg.JWTOptions.Key),
		Timeout:          cfg.JWTOptions.Timeout,
		MaxRefresh:       cfg.JWTOptions.MaxRefresh,
		Authenticator:    authenticator(),
		PayloadFunc:      payloadFunc(),
		Unauthorized:     unauthorized(),
		LoginResponse:    loginResponse(),
		RefreshResponse:  refreshResponse(),
		TokenLookup:      "header: Authorization, query: token",
		TimeFunc:         time.Now,
	})

	return jwtMiddleware
}

func authenticator() func(ctx *gin.Context) (any, error) {
	return func(c *gin.Context) (any, error) {
		var login loginInfo
		var user *model.Users
		var err error

		if err := c.ShouldBindBodyWith(&login, binding.JSON); err != nil {
			c.Set("VALIDATION_ERROR", validator.ParseValidationError(err))
			return nil, auth.ErrMissingLoginValues
		}

		err = middleware.WaterWall(c, "login:"+c.ClientIP())
		if err != nil {
			return nil, err
		}

		db := store.Client().DB()
		userStore := store.Client().User()

		switch {
		case verification.Phone(login.Username):
			user, err = userStore.Get(c, db, login.Username, options.WithQuery("phone = ?"))
		default:
			user, err = userStore.Get(c, db, login.Username, options.WithQuery("eid = ?"))
		}

		if err != nil {
			log.Errorf("get user information failed: %s", err.Error())
			return "", auth.ErrFailedAuthentication
		}

		if err := user.ComparePasswordHash(login.Password); err != nil {
			return "", auth.ErrFailedAuthentication
		}
		return user, nil
	}
}

func refreshResponse() func(c *gin.Context, _ int, token string, expire time.Time) {
	return func(c *gin.Context, _ int, token string, expire time.Time) {
		core.WriteResponse(c, tokenRsp{token, expire})
	}
}

func loginResponse() func(c *gin.Context, _ int, token string, expire time.Time) {
	return func(c *gin.Context, _ int, token string, expire time.Time) {
		core.WriteResponse(c, tokenRsp{token, expire})
	}
}

func payloadFunc() func(data any) auth.MapClaims {
	return func(data any) auth.MapClaims {
		claims := auth.MapClaims{
			"iss": APIServerIssuer,
			"aud": APIServerAudience,
		}
		if u, ok := data.(*model.Users); ok {
			claims["sub"] = u.EID
		}

		return claims
	}
}

func unauthorized() func(c *gin.Context, _ int, err error) {
	return func(c *gin.Context, _ int, err error) {
		var errCode int

		coder := errors.ParseCoder(err)
		if coder.Code() != 1 {
			if errors.IsCode(err, code.ErrNeedCaptcha) {
				cfg := config.GetConfigIns(nil).CaptchaOptions
				sign := idutil.GetInstanceID(idutil.GenUint64ID(), "")

				if err := middleware.SetCaptchaInfo(c, sign, cfg.AppID, cfg.AppSecretKey, "login:"+c.ClientIP(), time.Minute*10); err != nil {
					core.WriteResponse(c, nil, core.WithError(err), core.WithAbort())
					return
				}

				core.WriteResponse(
					c,
					middleware.NeedCaptchaRsp{
						Signature:    sign,
						CaptchaAppID: strconv.FormatUint(cfg.AppID, 10),
					},
					core.WithError(err), core.WithAbort(),
				)
				return
			}

			core.WriteResponse(c, nil, core.WithError(err), core.WithAbort())
			return
		}

		switch err {
		case auth.ErrFailedTokenCreation:
			errCode = code.ErrFailedTokenCreation
		case auth.ErrExpiredToken:
			errCode = code.ErrExpiredToken
		case auth.ErrFailedAuthentication:
			errCode = code.ErrFailedAuthentication

		case auth.ErrEmptyAuthHeader, auth.ErrEmptyParamToken, auth.ErrEmptyQueryToken:
			errCode = code.ErrEmptyToken

		case auth.ErrMissingExpField, auth.ErrWrongFormatOfExp, auth.ErrInvalidAuthHeader, auth.ErrInvalidSigningAlgorithm:
			log.Debugf("token 无效, err: %s", err.Error())
			errCode = code.ErrInvalidToken

		case auth.ErrMissingLoginValues:
			errCode = code.ErrValidation
			validationErr, exists := c.Get("VALIDATION_ERROR")
			if exists {
				core.WriteResponse(
					c,
					validationErr,
					core.WithError(errors.Code(code.ErrValidation, err.Error())),
					core.WithAbort(),
				)
				return
			}
		}

		core.WriteResponse(c, nil, core.WithError(errors.Code(errCode, err.Error())), core.WithAbort())
	}
}
