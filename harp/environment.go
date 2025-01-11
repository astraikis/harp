package main

var values = map[string]interface{}{}

func define(name string, value interface{}) {
	values[name] = value
}

func getValue(name token) interface{} {
	if val, ok := values[name.Lexeme]; ok {
		return val
	}

	return nil
}
