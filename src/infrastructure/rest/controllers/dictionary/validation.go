package dictionary

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gbrayhan/microservices-go/src/domain/constants"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	"github.com/go-playground/validator/v10"
)

func updateValidation(request map[string]any) error {
	var errorsValidation []string
	validationMap := map[string]string{
		"status":           "required,status_enum",
		"name":             "required,gt=3,lt=100",
		"type":             "required,gt=3,lt=100",
		"desc":             "omitempty,gt=3,lt=200",
		"is_generate_file": "omitempty",
	}

	validate := validator.New()

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

	err = validate.RegisterValidation("status_enum", func(fl validator.FieldLevel) bool {
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
