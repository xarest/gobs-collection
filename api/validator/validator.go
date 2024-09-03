package validator

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/xarest/gobs"

	v "github.com/go-playground/validator/v10"
)

type Validator struct {
	v *v.Validate
}

// Init implements gobs.IServiceInit.
func (v *Validator) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{
			&UserName{},
		},
	}, nil
}

func (vt *Validator) Setup(ctx context.Context, deps gobs.Dependencies) error {
	vt.v = v.New()
	for _, dep := range deps {
		if validator, ok := dep.(IValidator); ok {
			vFunc, tag := validator.Register()
			if err := vt.v.RegisterValidation(tag, vFunc); err != nil {
				return err
			}
		}
	}
	return nil
}

func (vt *Validator) Validate(i interface{}) error {
	if err := vt.v.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return err
	}
	return nil
}

var _ gobs.IServiceInit = (*Validator)(nil)
var _ gobs.IServiceSetup = (*Validator)(nil)
var _ echo.Validator = (*Validator)(nil)
