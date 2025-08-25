// ValidatorWrapper.go
package controllers

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/gbrayhan/microservices-go/src/domain/constants"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	"github.com/go-playground/validator/v10"
)

type CommonValidator struct {
	validate       *validator.Validate
	validationMap  map[string]string
	errorsMessages []string
}

// NewCommonValidator 创建新的通用验证器实例
func NewCommonValidator(validationMap map[string]string) *CommonValidator {
	v := &CommonValidator{
		validate:       validator.New(),
		validationMap:  validationMap,
		errorsMessages: make([]string, 0),
	}

	v.registerCustomValidations()
	return v
}

func (v *CommonValidator) registerCustomValidations() {
	_ = v.validate.RegisterValidation("update_validation", func(fl validator.FieldLevel) bool {
		m, ok := fl.Field().Interface().(map[string]any)
		if !ok {
			return false
		}

		for k, rule := range v.validationMap {
			if val, exists := m[k]; exists {
				errValidate := v.validate.Var(val, rule)
				if errValidate != nil {
					validatorErr := errValidate.(validator.ValidationErrors)
					v.errorsMessages = append(
						v.errorsMessages,
						fmt.Sprintf("%s does not satisfy condition %v=%v", k, validatorErr[0].Tag(), validatorErr[0].Param()),
					)
				}
			}
		}
		return true
	})

	_ = v.validate.RegisterValidation("status_enum", func(fl validator.FieldLevel) bool {
		value := fl.Field().Interface()
		switch v := value.(type) {
		case float64:
			// 转换为字符串并检查是否匹配枚举值
			strValue := fmt.Sprintf("%.0f", v)
			return strValue == constants.StatusEnable || strValue == constants.StatusDisable
		case int:
			strValue := fmt.Sprintf("%d", v)
			return strValue == constants.StatusEnable || strValue == constants.StatusDisable
		case string:
			return v == constants.StatusEnable || v == constants.StatusDisable
		default:
			return false
		}
	})

	_ = v.validate.RegisterValidation("custom_phone", func(fl validator.FieldLevel) bool {
		phone, ok := fl.Field().Interface().(string)
		if !ok || phone == "" {
			return true // 允许空值由 omitempty 处理
		}
		// 自定义手机号正则表达式（示例为中国手机号）
		match := regexp.MustCompile(`^\+?\d{10,15}$`).MatchString(phone)
		return match
	})

}

func (v *CommonValidator) ValidateUpdate(request map[string]any) error {
	// 重置错误消息
	v.errorsMessages = make([]string, 0)

	// 基本空值检查
	for k, val := range request {
		if val == "" {
			v.errorsMessages = append(v.errorsMessages, fmt.Sprintf("%s cannot be empty", k))
		}
	}

	// 执行自定义验证
	err := v.validate.Var(request, "update_validation")
	if err != nil {
		return domainErrors.NewAppError(err, domainErrors.UnknownError)
	}

	if len(v.errorsMessages) > 0 {
		return domainErrors.NewAppError(errors.New(strings.Join(v.errorsMessages, ", ")), domainErrors.ValidationError)
	}

	return nil
}
