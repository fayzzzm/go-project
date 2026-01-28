package validate

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fayzzzm/go-project/internal/domain"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	defaultValidator = NewValidator()
)

// RegisterGinValidator replaces the default Gin validator with our custom one.
func RegisterGinValidator() {
	binding.Validator = &ginValidator{v: defaultValidator.v}
}

type ginValidator struct {
	v *validator.Validate
}

func (gv *ginValidator) ValidateStruct(obj interface{}) error {
	if reflect.Indirect(reflect.ValueOf(obj)).Kind() != reflect.Struct {
		return nil
	}
	return gv.v.Struct(obj)
}

func (gv *ginValidator) Engine() interface{} {
	return gv.v
}

type ValidationErrorDetail struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

// Validate is a package-level helper using the default validator.
func Validate(s interface{}) error {
	return defaultValidator.Validate(s)
}

// GetErrors is a package-level helper using the default validator.
func GetErrors(err error) []ValidationErrorDetail {
	return defaultValidator.GetErrors(err)
}

type Validator struct {
	v *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()

	v.SetTagName("validate")

	// Use snake_case for field names in errors
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validator{v: v}
}

// Validate validates a struct and returns a friendly error if it fails.
func (vr *Validator) Validate(s interface{}) error {
	return vr.v.Struct(s)
}

// GetErrors returns a slice of ValidationErrorDetail for the given error.
func (vr *Validator) GetErrors(err error) []ValidationErrorDetail {
	var details []ValidationErrorDetail
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			details = append(details, ValidationErrorDetail{
				Field:  fe.Field(),
				Reason: vr.msgForTag(fe),
			})
		}
	}
	return details
}

func (vr *Validator) msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case domain.ValidationTagRequired:
		return "This field is required"
	case domain.ValidationTagEmail:
		return "Invalid email format"
	case "min":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("Must be at least %s characters long", fe.Param())
		}
		return fmt.Sprintf("Must be at least %s", fe.Param())
	case "max":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("Must be no more than %s characters long", fe.Param())
		}
		return fmt.Sprintf("Must be no more than %s", fe.Param())
	case domain.ValidationTagUUID:
		return "Must be a valid UUID"
	default:
		return fmt.Sprintf("Failed validation on tag '%s'", fe.Tag())
	}
}
