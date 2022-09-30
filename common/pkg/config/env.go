package config

import (
	"errors"
	"os"
)

//use export ENVIRONMENT=testing set global environment
const (
	EnvTesting    = Environment("testing")
	EnvProduction = Environment("production")
)

type Environment string

func (env *Environment) String() string {
	return string(*env)
}

func (env *Environment) Production() Environment {
	return EnvProduction
}

func (env *Environment) Testing() Environment {
	return EnvTesting
}

func (env Environment) Invalid() bool {
	return env != EnvTesting && env != EnvProduction
}

// NewGlobalEnvironment 读取系统全局配置的环境变量
func NewGlobalEnvironment() (Environment, error) {
	environment, ok := os.LookupEnv("ENVIRONMENT")
	if !ok {
		return "", errors.New("system environment:ENVIRONMENT not found")
	}

	env := Environment(environment)
	if env != EnvTesting && env != EnvProduction {
		return "", errors.New("environment not support, must be production, testing, development")
	}

	return env, nil
}
