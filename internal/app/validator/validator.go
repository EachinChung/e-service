package validator

import (
	"context"
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
	"github.com/eachinchung/e-service/internal/app/store/casbin"
)

var Trans ut.Translator

const (
	password  = "password"
	phone     = "phone"
	eid       = "eid"
	isNotRole = "is_not_role"
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
		if err := v.RegisterValidation(eid, eidValidation); err != nil {
			return err
		}
		if err := v.RegisterValidation(isNotRole, isNotRoleValidation); err != nil {
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
			registerTranslator(phone, "必须为合法的中国大陆手机号码"),
			translate,
		); err != nil {
			return err
		}
		if err := v.RegisterTranslation(
			password,
			Trans,
			registerTranslator(password, "密码必须存在特殊字符、大小写字母和数字"),
			translate,
		); err != nil {
			return err
		}
		if err := v.RegisterTranslation(
			eid,
			Trans,
			registerTranslator(eid, "用户名不能以数字开头，可以使用6-20位字母、数字、下划线或减号组合而成"),
			translate,
		); err != nil {
			return err
		}
		if err := v.RegisterTranslation(
			isNotRole,
			Trans,
			registerTranslator(isNotRole, "用户名不能为敏感名称"),
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

// eidValidation eid 校验
func eidValidation(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	rgx := regexp.MustCompile(`^[a-zA-Z][a-zA-Z\d_-]{5,19}$`)
	return rgx.MatchString(val)
}

// isNotRoleValidation 角色名校验，用户名不能为 rbac 的角色名
func isNotRoleValidation(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	permissions := casbin.GetPermissionsForUser(context.Background(), val)
	return len(permissions) == 0
}
