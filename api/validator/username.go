package validator

import (
	"regexp"

	v "github.com/go-playground/validator/v10"
)

var UserNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,64}$`)

type UserName struct {
	// log logger.ILogger
}

// Uncomment these snippets if you want to use the logger service

// func (uv *UserName) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
// 	return &gobs.ServiceLifeCycle{
// 		Deps: gobs.Dependencies{
// 			&logger.ILogger{},
// 		},
// 	}, nil
// }

// func (uv *UserName) Setup(ctx context.Context, deps gobs.Dependencies) error {
// 	deps.Assign(&uv.log)
// }

func (uv *UserName) Register() (v.Func, string) {
	return func(fl v.FieldLevel) bool {
		return UserNameRegex.MatchString(fl.Field().String())
	}, "username"
}

var _ IValidator = (*UserName)(nil)
