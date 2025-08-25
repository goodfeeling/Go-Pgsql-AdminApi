package dictionary

import "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"

var customRules = map[string]string{
	"status":           "required,status_enum",
	"name":             "required,gt=3,lt=100",
	"type":             "required,gt=3,lt=100",
	"desc":             "omitempty,gt=3,lt=200",
	"is_generate_file": "omitempty",
}

func updateValidation(request map[string]any) error {
	validator := controllers.NewCommonValidator(customRules)
	return validator.ValidateUpdate(request)
}
