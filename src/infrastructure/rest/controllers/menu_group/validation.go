package menu_group

import "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"

var customRules = map[string]string{
	"name":   "required|max:255",
	"path":   "required|max:255",
	"sort":   "required|numeric",
	"status": "required|status_enum",
}

func updateValidation(request map[string]any) error {
	validator := controllers.NewCommonValidator(customRules)
	return validator.ValidateUpdate(request)
}
