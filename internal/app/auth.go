package app

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/eachinchung/component-base/core"
	"github.com/eachinchung/component-base/db/options"
	"github.com/eachinchung/component-base/middleware/auth"
	"github.com/eachinchung/component-base/verification"
	"github.com/eachinchung/errors"
	"github.com/eachinchung/log"

	"github.com/eachinchung/e-service/internal/app/config"
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
	Username string `form:"username" json:"username" binding:"required,max=64"`
	Password string `form:"password" json:"password" binding:"required,min=6,password"`
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

func authenticator() func(ctx *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var user *model.Users
		var err error

		login, err := parseWithBody(c)
		if err != nil {
			return "", auth.ErrMissingLoginValues
		}

		db := store.Client().DB()
		userStore := store.Client().User()

		switch {
		case verification.Phone(login.Username):
			user, err = userStore.Get(c, db, login.Username, options.WithQuery("phone = ?"))
		default:
			user, err = userStore.Get(c, db, login.Username, options.WithQuery("username = ?"))
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

func parseWithBody(c *gin.Context) (loginInfo, error) {
	var login loginInfo
	if err := c.ShouldBindJSON(&login); err != nil {
		return loginInfo{}, auth.ErrMissingLoginValues
	}

	return login, nil
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

func payloadFunc() func(data interface{}) auth.MapClaims {
	return func(data interface{}) auth.MapClaims {
		claims := auth.MapClaims{
			"iss": APIServerIssuer,
			"aud": APIServerAudience,
		}
		if u, ok := data.(*model.Users); ok {
			claims["sub"] = u.Username
		}

		return claims
	}
}

func unauthorized() func(c *gin.Context, _ int, err error) {
	return func(c *gin.Context, _ int, err error) {
		c.Abort()

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
		case auth.ErrEmptyAuthHeader, auth.ErrEmptyParamToken, auth.ErrEmptyQueryToken:
			errCode = code.ErrEmptyToken
		}

		core.WriteResponse(c, nil, core.WithError(errors.Code(errCode, err.Error())))
	}
}
