package utilities

import (
	"mime/multipart"
	"strconv"

	"github.com/go-playground/validator/v10"
)

// IValidator helper with validator
type IValidator interface {
	Integer(input string, tag string) (int, error)
	Integer64(input string, tag string) (int64, error)
	String(input string, tag string) (string, error)
	Float(input string, tag string) (float64, error)
	Bool(input string, tag string) (bool, error)
	File(input *multipart.FileHeader, tag string) (*multipart.FileHeader, error)
}

// ValidatorImpl helper
type ValidatorImpl struct {
	*validator.Validate
}

// New validator utility
func New() ValidatorImpl {
	validate := validator.New()
	return ValidatorImpl{
		Validate: validate,
	}
}

// Integer get input string, check valid and return value
func (validator ValidatorImpl) Integer(input string, tag string) (int, error) {
	err := validator.Validate.Var(input, tag)
	if err != nil {
		return 0, err
	}
	number, _ := strconv.Atoi(input)
	return number, nil
}

// Integer64 get input string, check valid and return value
func (validator ValidatorImpl) Integer64(input string, tag string) (int64, error) {
	err := validator.Validate.Var(input, tag)
	if err != nil {
		return 0, err
	}
	number, _ := strconv.ParseInt(input, 10, 64)
	return number, nil
}

// String get input string, check valid and return value
func (validator ValidatorImpl) String(input string, tag string) (string, error) {
	err := validator.Validate.Var(input, tag)
	if err != nil {
		return "", err
	}
	return input, nil
}

// Float get input string, check valid and return value
func (validator ValidatorImpl) Float(input string, tag string) (float64, error) {
	err := validator.Validate.Var(input, tag)
	if err != nil {
		return 0, err
	}
	number, _ := strconv.ParseFloat(input, 64)
	return number, nil
}

// Bool get input string, check valid and return value
func (validator ValidatorImpl) Bool(input string, tag string) (bool, error) {
	err := validator.Validate.Var(input, tag)
	if err != nil {
		return false, err
	}
	value, _ := strconv.ParseBool(input)
	return value, nil
}

// File get input file, check valid and return value
func (validator ValidatorImpl) File(input *multipart.FileHeader, tag string) (*multipart.FileHeader, error) {
	err := validator.Validate.Var(input, tag)
	if err != nil {
		return nil, err
	}
	return input, nil
}
