package menu_group

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	"github.com/go-playground/validator/v10"
)

func updateValidation(request map[string]any) error {
	var errorsValidation []string
	for k, v := range request {
		if v == "" {
			errorsValidation = append(errorsValidation, fmt.Sprintf("%s cannot be empty", k))
		}
	}

	validationMap := map[string]string{
		"user_name": "omitempty,gt=3,lt=100",
		"email":     "omitempty,email",
		"phone":     "omitempty,custom_phone",
		"nick_name": "omitempty",
	}

	validate := validator.New()

	// 注册自定义电话号码验证规则
	_ = validate.RegisterValidation("custom_phone", func(fl validator.FieldLevel) bool {
		phone, ok := fl.Field().Interface().(string)
		if !ok || phone == "" {
			return true // 允许空值由 omitempty 处理
		}
		// 自定义手机号正则表达式（示例为中国手机号）
		match := regexp.MustCompile(`^\+?\d{10,15}$`).MatchString(phone)
		return match
	})

	// 保留原有的 update_validation 逻辑
	err := validate.RegisterValidation("update_validation", func(fl validator.FieldLevel) bool {
		m, ok := fl.Field().Interface().(map[string]any)
		if !ok {
			return false
		}
		for k, rule := range validationMap {
			if val, exists := m[k]; exists {
				errValidate := validate.Var(val, rule)
				if errValidate != nil {
					validatorErr := errValidate.(validator.ValidationErrors)
					errorsValidation = append(
						errorsValidation,
						fmt.Sprintf("%s does not satisfy condition %v=%v", k, validatorErr[0].Tag(), validatorErr[0].Param()),
					)
				}
			}
		}
		return true
	})
	if err != nil {
		return domainErrors.NewAppError(err, domainErrors.UnknownError)
	}

	err = validate.Var(request, "update_validation")
	if err != nil {
		return domainErrors.NewAppError(err, domainErrors.UnknownError)
	}
	if len(errorsValidation) > 0 {
		return domainErrors.NewAppError(errors.New(strings.Join(errorsValidation, ", ")), domainErrors.ValidationError)
	}
	return nil

}
