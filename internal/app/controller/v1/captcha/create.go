package captcha

import (
	"fmt"

	captcha "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/captcha/v20190722"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

func qwe() {
	// 实例化一个认证对象，入参需要传入腾讯云账户secretId，secretKey,此处还需注意密钥对的保密
	// 密钥可前往https://console.cloud.tencent.com/cam/capi网站进行获取
	credential := common.NewCredential(
		"SecretId",
		"SecretKey",
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "captcha.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := captcha.NewClient(credential, "", cpf)

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := captcha.NewDescribeCaptchaResultRequest()

	request.CaptchaType = common.Uint64Ptr(9)
	request.Ticket = common.StringPtr("75675")
	request.UserIp = common.StringPtr("127,0,0,1")
	request.Randstr = common.StringPtr("ljkl")
	request.CaptchaAppId = common.Uint64Ptr(1224324)
	request.AppSecretKey = common.StringPtr("jhkhkhk")

	// 返回的resp是一个DescribeCaptchaResultResponse的实例，与请求对象对应
	response, err := client.DescribeCaptchaResult(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	if err != nil {
		panic(err)
	}
	// 输出json格式的字符串回包
	fmt.Printf("%s", response.ToJsonString())
}
