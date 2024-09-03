package validator

import v "github.com/go-playground/validator/v10"

type IValidator interface {
	Register() (v.Func, string)
}
