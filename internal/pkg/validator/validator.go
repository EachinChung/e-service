package validator

import (
	"errors"
	"reflect"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"

	"github.com/eachinchung/component-base/verification"
)

var Trans ut.Translator

const (
	password = "password"
	phone    = "phone"
	username = "username"
)

func InitValidator() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		if err := v.RegisterValidation(phone, phoneValidation); err != nil {
			return err
		}
		if err := v.RegisterValidation(password, passwordValidation); err != nil {
			return err
		}
		if err := v.RegisterValidation(username, usernameValidation); err != nil {
			return err
		}

		zhT := zh.New()
		uni := ut.New(zhT, zhT)
		Trans, ok = uni.GetTranslator("zh")
		if !ok {
			return errors.New("uni.GetTranslator failed")
		}

		if err := zhTranslations.RegisterDefaultTranslations(v, Trans); err != nil {
			return err
		}

		if err := v.RegisterTranslation(
			phone,
			Trans,
			registerTranslator(phone, "{0}必须为合法的中国大陆手机号码"),
			translate,
		); err != nil {
			return err
		}
		if err := v.RegisterTranslation(
			password,
			Trans,
			registerTranslator(password, "{0}必须存在特殊字符、大小写字母和数字"),
			translate,
		); err != nil {
			return err
		}
		if err := v.RegisterTranslation(
			username,
			Trans,
			registerTranslator(username, "{0}不能以数字开头，可以使用6-20位字母、数字、下划线或减号组合而成"),
			translate,
		); err != nil {
			return err
		}
	}
	return nil
}

// ParseValidationError 解析错误信息
func ParseValidationError(err error) map[string]string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		return removeTopStruct(errs.Translate(Trans))
	}

	return nil
}

// removeTopStruct 去除提示信息中的结构体名称
func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}

// registerTranslator 为自定义字段添加翻译功能
func registerTranslator(tag string, msg string) validator.RegisterTranslationsFunc {
	return func(trans ut.Translator) error {
		if err := trans.Add(tag, msg, false); err != nil {
			return err
		}
		return nil
	}
}

// translate 自定义字段的翻译方法
func translate(trans ut.Translator, fe validator.FieldError) string {
	msg, err := trans.T(fe.Tag(), fe.Field())
	if err != nil {
		panic(fe.(error).Error())
	}
	return msg
}

// passwordValidation 密码校验
func passwordValidation(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	return verification.PasswordPower(val)
}

// phoneValidation 手机号校验
func phoneValidation(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	return verification.Phone(val)
}

// usernameValidation 用户名校验
func usernameValidation(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	rgx := regexp.MustCompile(`^[a-zA-Z][a-zA-Z\d_-]{5,19}$`)
	return rgx.MatchString(val)
}
