package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// validate Indian mobile numbers
var indianPhoneRegex = regexp.MustCompile(`^[6-9]\d{9}$`)

func IndianPhoneValidator(fl validator.FieldLevel) bool {
	return indianPhoneRegex.MatchString(fl.Field().String())
}
