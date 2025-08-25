package scheduled_task

import "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"

var customRules = map[string]string{
	"task_name":        "required|max:255",
	"task_description": "required",
	"cron_expression":  "required|max:255",
	"exec_type":        "required|max:50",
	"task_type":        "required|max:100",
	"task_params":      "required",
}

func updateValidation(request map[string]any) error {
	validator := controllers.NewCommonValidator(customRules)
	return validator.ValidateUpdate(request)
}
