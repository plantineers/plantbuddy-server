package plant

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func InitializeValidator() {
	validate = validator.New()
}
