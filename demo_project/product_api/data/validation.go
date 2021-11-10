package data

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator"
)

type ValidationError struct {
	validator.FieldError
}

// generalized function which returns an error
func (v ValidationError) Error() string {
	return fmt.Sprintf(
		"%s %s %s",
		v.Namespace(),
		v.Field(),
		v.Tag(),
	)
}

// valerrors is a collection of ValidationError struct
type valerrors []ValidationError

// Errors converts the slice into a string slice
func (v valerrors) Errors() []string {
	errs := []string{}
	for _, err := range v {
		errs = append(errs, err.Error())
	}

	return errs
}

type Validation struct {
	validate *validator.Validate
}

func NewValidation() *Validation {
	v := validator.New()
	v.RegisterValidation("sku", validateSKU)

	return &Validation{validate: v}
}

func (v *Validation) Validate(i interface{}) valerrors {
	errs := v.validate.Struct(i).(validator.ValidationErrors) // assertion
	if len(errs) == 0 {
		return nil
	}

	// this block executes if the above if block doesnt satisfy the condition
	var returnerr []ValidationError
	// looping over the slice of errors received by the above statement after running the validation
	for _, err := range errs {
		ve := ValidationError{err}
		returnerr = append(returnerr, ve)
	}

	return returnerr
}

func validateSKU(xd validator.FieldLevel) bool {
	// SKU must be in the format abc-abc-abc
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`) // specify the regex here
	sku := re.FindAllString(xd.Field().String(), -1) // finds the strings that matches the regexp

	if len(sku) == 1 {
		return true // returns true when kyu is just a single string
	}

	return false
}
