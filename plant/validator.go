package plant

import "github.com/go-playground/validator/v10"

func InitializeValidator() {
	validate = validator.New()
}
