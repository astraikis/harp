package interpreter

import (
	"github.com/astraikis/harp/internal/models"
)

type environment struct {
	values map[string]interface{}
	parent *environment
}

func defineValue(name string, value interface{}, currentEnvironment *environment) {
	currentEnvironment.values[name] = value
}

func getValue(name models.Token, currentEnvironment *environment) interface{} {
	if val, ok := currentEnvironment.values[name.Lexeme]; ok {
		return val
	}

	if currentEnvironment.parent != nil {
		return getValue(name, currentEnvironment.parent)
	}

	return nil

}

func assignValue(name models.Token, value interface{}, currentEnvironment *environment) {
	if _, ok := currentEnvironment.values[name.Lexeme]; ok {
		currentEnvironment.values[name.Lexeme] = value
		return
	}

	if currentEnvironment.parent != nil {
		assignValue(name, value, currentEnvironment.parent)
		return
	}

	panic("Undefined variable.")
}
