package interpreter

import "github.com/astraikis/harp/internal/models"

type environment struct {
	values map[string]interface{}
	parent *environment
}

func defineValue(name string, value interface{}, currEnvironment environment) {
	currEnvironment.values[name] = value
}

func getValue(name models.Token, currEnvironment environment) interface{} {
	if val, ok := currEnvironment.values[name.Lexeme]; ok {
		return val
	}

	if currEnvironment.parent != nil {
		return getValue(name, *currEnvironment.parent)
	}

	return nil
}

func assignValue(name models.Token, value interface{}, currEnvironment environment) {
	if _, ok := currEnvironment.values[name.Lexeme]; ok {
		currEnvironment.values[name.Lexeme] = value
		return
	}

	if currEnvironment.parent != nil {
		assignValue(name, value, *currEnvironment.parent)
	}

	panic("Undefined variable.")
}
