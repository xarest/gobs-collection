package utils

import (
	"context"
	"reflect"

	"github.com/xarest/gobs"
)

func NewGobsInstannce[T gobs.IServiceSetup](deps ...gobs.IService) T {
	var instance T
	instanceType := reflect.TypeOf(instance).Elem()
	instanceValue := reflect.New(instanceType).Interface().(T)
	if err := instanceValue.Setup(context.Background(), deps); err != nil {
		panic(err)
	}
	return instanceValue
}
