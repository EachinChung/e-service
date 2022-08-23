package middleware

import (
	"context"
	"time"

	"github.com/eachinchung/errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/eachinchung/e-service/internal/app/config"
	"github.com/eachinchung/e-service/internal/app/storage"
	"github.com/eachinchung/e-service/internal/pkg/captcha"
	"github.com/eachinchung/e-service/internal/pkg/code"
)

const (
	redisPrefix     = "water-wall:"
	redisInfoPrefix = "water-wall:info:"
)

type NeedCaptchaRsp struct {
	Signature    string `json:"signature"`
	CaptchaAppID string `json:"captcha_app_id"`
}

type waterWallInfo struct {
	Ticket    string `form:"ticket" json:"ticket" binding:"required"`
	RandStr   string `form:"rand_str" json:"rand_str" binding:"required"`
	Signature string `form:"signature" json:"signature" binding:"required"`
}

type captchaInfo struct {
	CaptchaAppID        uint64 `redis:"captcha_app_id"`
	CaptchaAppSecretKey string `redis:"captcha_app_secret_key"`
	Identifier          string `redis:"identifier"`
}

type waterWallConfig struct {
	burst  int64
	period time.Duration
}

type Option func(*waterWallConfig)

func WithLimit(burst int64, period time.Duration) Option {
	return func(cfg *waterWallConfig) {
		cfg.burst = burst
		cfg.period = period
	}
}

func SetCaptchaInfo(
	ctx context.Context,
	sign string,
	appID uint64,
	appSecretKey, identifier string,
	expiration time.Duration,
) error {
	rdb := storage.Client()
	cInfoKey := redisInfoPrefix + sign
	err := rdb.HSetAllWithExpire(
		ctx,
		cInfoKey,
		&captchaInfo{
			CaptchaAppID:        appID,
			CaptchaAppSecretKey: appSecretKey,
			Identifier:          identifier,
		},
		expiration,
	)

	return errors.WithCode(err, code.ErrCaptchaBusy, "water wall: redis HSetAll error")
}

func WaterWall(c *gin.Context, identifier string, opts ...Option) error {
	wc := &waterWallConfig{}
	for _, opt := range opts {
		opt(wc)
	}

	rdb := storage.Client()

	if wc.burst > 0 && wc.period > 0 {
		cKey := redisPrefix + identifier
		v, err := rdb.Incr(c, cKey)
		if err != nil {
			return errors.WithCode(err, code.ErrCaptchaBusy, "water wall: redis INCR error")
		}
		_ = rdb.Expire(c, cKey, wc.period)

		if v < wc.burst {
			return nil
		}
	}

	cfg := config.GetConfigIns(nil).CaptchaOptions
	var info waterWallInfo
	if err := c.ShouldBindBodyWith(&info, binding.JSON); err != nil {
		return errors.WithCodef(err, code.ErrNeedCaptcha, err.Error())
	}

	cInfoKey := redisInfoPrefix + info.Signature
	var captchaInfo captchaInfo
	if err := rdb.HGetAll(c, cInfoKey, &captchaInfo); err != nil {
		return errors.WithCode(err, code.ErrCaptchaBusy, "water wall: redis HGETALL error")
	}

	if identifier != captchaInfo.Identifier {
		return errors.Code(code.ErrNetworkUnsafe, "water wall: identifier error")
	}

	if !cfg.IsVerify {
		return nil
	}

	client := captcha.GetClientOr(&captcha.ClientConfig{
		SecretID:  cfg.SecretID,
		SecretKey: cfg.SecretKey,
	})

	if err := captcha.New(
		client,
		captchaInfo.CaptchaAppID,
		captchaInfo.CaptchaAppSecretKey,
	).Verify(c, info.Ticket, info.RandStr, c.ClientIP()); err != nil {
		return err
	}
	return nil
}
