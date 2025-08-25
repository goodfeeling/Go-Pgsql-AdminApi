package menu

import "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"

var customRules = map[string]string{
	"component":  "required|max:191",
	"title":      "required|max:191",
	"name":       "required|max:191",
	"path":       "required|max:191",
	"hidden":     "required",
	"keep_alive": "required",
	"parent_id":  "required|max:11",
	"icon":       "required|max:191",
	"sort":       "required",
}

func updateValidation(request map[string]any) error {
	validator := controllers.NewCommonValidator(customRules)
	return validator.ValidateUpdate(request)
}
