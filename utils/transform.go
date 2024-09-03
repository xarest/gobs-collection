package utils

import (
	"slices"

	"github.com/mitchellh/mapstructure"
)

func ToMap[K comparable, T any](slice []T, key func(T) K, init map[K]T) map[K]T {
	var m map[K]T
	if init != nil {
		m = init
	} else {
		m = make(map[K]T)
	}
	for _, entry := range slice {
		m[key(entry)] = entry
	}
	return m
}

func GetListFields(params any) ([]string, error) {
	mParams := map[string]any{}
	if err := mapstructure.Decode(params, mParams); err != nil {
		return nil, err
	}
	return slices.Collect(func(yield func(col string) bool) {
		for k := range mParams {
			if !yield(k) {
				return
			}
		}
	}), nil
}
