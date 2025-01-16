package interpreter

import "github.com/astraikis/harp/internal/models"

var values = map[string]interface{}{}

func defineValue(name string, value interface{}) {
	values[name] = value
}

func getValue(name models.Token) interface{} {
	if val, ok := values[name.Lexeme]; ok {
		return val
	}

	return nil
}

func assignValue(name models.Token, value interface{}) {
	if _, ok := values[name.Lexeme]; ok {
		values[name.Lexeme] = value
		return
	}

	panic("Undefined variable.")
}
