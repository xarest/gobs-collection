package config

type IConfiguration interface {
	Parse(result interface{}) error
}

func NewIConfig() IConfiguration {
	return &EnvConfig{}
}
