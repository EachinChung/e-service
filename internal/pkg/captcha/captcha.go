package captcha

//goland:noinspection SpellCheckingInspection
import (
	"context"
	"sync"

	"github.com/eachinchung/errors"
	"github.com/eachinchung/log"
	captcha "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/captcha/v20190722"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tencenterrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"

	"github.com/eachinchung/e-service/internal/pkg/code"
)

type Captcha interface {
	Verify(ctx context.Context, ticket, randStr, userIP string) error
}

type ClientConfig struct {
	SecretID, SecretKey string
}

type application struct {
	client      *captcha.Client
	id          uint64
	secretKey   string
	captchaType uint64
}

var once sync.Once

func GetClientOr(cfg *ClientConfig) (client *captcha.Client) {
	once.Do(func() {
		credential := common.NewCredential(cfg.SecretID, cfg.SecretKey)
		cpf := profile.NewClientProfile()
		cpf.HttpProfile.Endpoint = "captcha.tencentcloudapi.com"
		client, _ = captcha.NewClient(credential, "", cpf)
	})

	return client
}

func New(client *captcha.Client, id uint64, secretKey string) Captcha {
	return &application{
		client:      client,
		id:          id,
		secretKey:   secretKey,
		captchaType: 9,
	}
}

func (app *application) Verify(ctx context.Context, ticket, randStr, userIP string) error {
	req := captcha.NewDescribeCaptchaResultRequest()
	req.Ticket = &ticket
	req.Randstr = &randStr
	req.AppSecretKey = &app.secretKey
	req.CaptchaAppId = &app.id
	req.CaptchaType = &app.captchaType
	req.UserIp = &userIP
	resp, err := app.client.DescribeCaptchaResult(req)
	if e, ok := err.(*tencenterrors.TencentCloudSDKError); ok {
		log.L(ctx).Errorf("captcha verify failed: %s", e.Message)
		return errors.WithCodef(err, code.ErrCaptchaBusy, "验证码服务异常")
	}
	if err != nil {
		log.L(ctx).Errorf("captcha verify failed: %+v", err)
		return errors.WithCodef(err, code.ErrCaptchaBusy, "验证码服务异常")
	}

	log.L(ctx).Infof("verify captcha result: %+v", resp)
	if *resp.Response.CaptchaCode != 1 {
		return errors.Code(code.ErrCaptchaVerifyFailed, *resp.Response.CaptchaMsg)
	}

	return nil
}
