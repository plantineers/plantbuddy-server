package plant

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

// InitializeValidator initializes the validator for the plant package.
func InitializeValidator() {
	validate = validator.New()
}
