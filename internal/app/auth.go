package app

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/eachinchung/e-service/internal/app/validator"

	"github.com/gin-gonic/gin"

	"github.com/eachinchung/component-base/core"
	"github.com/eachinchung/component-base/db/options"
	"github.com/eachinchung/component-base/middleware/auth"
	"github.com/eachinchung/component-base/utils/idutil"
	"github.com/eachinchung/component-base/verification"
	"github.com/eachinchung/errors"
	"github.com/eachinchung/log"

	"github.com/eachinchung/e-service/internal/app/config"
	"github.com/eachinchung/e-service/internal/app/storage"
	"github.com/eachinchung/e-service/internal/app/store"
	"github.com/eachinchung/e-service/internal/app/store/model"
	"github.com/eachinchung/e-service/internal/pkg/code"
)

const (
	// APIServerAudience defines the value of jwt audience field.
	APIServerAudience = "api.service.eachin-life.com"

	// APIServerIssuer defines the value of jwt issuer field.
	APIServerIssuer = "e-service"
)

type loginInfo struct {
	Username  string  `form:"username" json:"username" binding:"required,max=20"`
	Password  string  `form:"password" json:"password" binding:"required"`
	Ticket    *string `form:"ticket" json:"ticket" binding:"omitempty"`
	RandStr   *string `form:"rand_str" json:"rand_str" binding:"omitempty"`
	Signature *string `form:"signature" json:"signature" binding:"omitempty"`
}

type needCaptchaRsp struct {
	Signature    string `json:"signature"`
	CaptchaAppID string `json:"captcha_app_id"`
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

		if err := c.ShouldBindJSON(&login); err != nil {
			c.Set("VALIDATION_ERROR", validator.ParseValidationError(err))
			return nil, auth.ErrMissingLoginValues
		}

		rdb := storage.Client()
		cKey := fmt.Sprintf(storage.KeyLoginIP, c.ClientIP())

		numberOfVisits, err := rdb.Incr(c, cKey)
		log.Errorf("%s: %s", cKey, numberOfVisits)
		if err == nil && numberOfVisits > 10 {
			_ = rdb.Expire(c, cKey, time.Hour)

			if login.RandStr == nil || login.Ticket == nil || login.Signature == nil {
				c.Set("NEED_CAPTCHA", login.Username)
				return nil, errors.Code(code.ErrNeedCaptcha, "需要验证码")
			}

			// TODO: check captcha
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
		c.JSON(http.StatusOK, gin.H{
			"token":  token,
			"expire": expire.Format(time.RFC3339),
		})
	}
}

func loginResponse() func(c *gin.Context, _ int, token string, expire time.Time) {
	return func(c *gin.Context, _ int, token string, expire time.Time) {
		c.JSON(http.StatusOK, gin.H{
			"token":  token,
			"expire": expire.Format(time.RFC3339),
		})
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

		switch err {
		case auth.ErrFailedTokenCreation:
			errCode = code.ErrFailedTokenCreation
		case auth.ErrExpiredToken:
			errCode = code.ErrExpiredToken
		case auth.ErrMissingExpField:
			errCode = code.ErrMissingExpField
		case auth.ErrWrongFormatOfExp:
			errCode = code.ErrWrongFormatOfExp
		case auth.ErrInvalidAuthHeader:
			errCode = code.ErrInvalidAuthHeader
		case auth.ErrInvalidSigningAlgorithm:
			errCode = code.ErrInvalidSigningAlgorithm
		case auth.ErrFailedAuthentication:
			errCode = code.ErrFailedAuthentication
		case auth.ErrMissingLoginValues:
			errCode = code.ErrValidation
			validationErr, exists := c.Get("VALIDATION_ERROR")
			if exists {
				core.WriteResponse(
					c,
					validationErr,
					core.WithError(errors.Code(code.ErrValidation, err.Error())),
				)
				return
			}
		case auth.ErrEmptyAuthHeader, auth.ErrEmptyParamToken, auth.ErrEmptyQueryToken:
			errCode = code.ErrEmptyToken
		}

		if errors.IsCode(err, code.ErrNeedCaptcha) {
			cfg := config.GetConfigIns(nil)
			sign := idutil.GetInstanceID(idutil.GenUint64ID(), "")

			rdb := storage.Client()
			cKey := fmt.Sprintf(storage.KeyLoginSign, sign)
			username, _ := c.Get("NEED_CAPTCHA")
			if err := rdb.Set(c, cKey, username, time.Minute*5); err != nil {
				core.WriteResponse(c, nil, core.WithError(errors.Code(code.ErrDatabase, err.Error())), core.WithAbort())
				return
			}

			core.WriteResponse(
				c,
				needCaptchaRsp{
					Signature:    sign,
					CaptchaAppID: strconv.FormatUint(cfg.TencentCloudOptions.CaptchaAppID, 10),
				},
				core.WithError(err), core.WithAbort(),
			)
			return
		}
		core.WriteResponse(c, nil, core.WithError(errors.Code(errCode, err.Error())), core.WithAbort())
	}
}
