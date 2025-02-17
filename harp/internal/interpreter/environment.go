package interpreter

type Environment struct {
	values map[string]interface{}
	parent *Environment
}

func DefineValue(name string, value interface{}, currentEnvironment *Environment) {
	currentEnvironment.values[name] = value
}

func GetValue(name string, currentEnvironment *Environment) interface{} {
	if val, ok := currentEnvironment.values[name]; ok {
		return val
	}

	if currentEnvironment.parent != nil {
		return GetValue(name, currentEnvironment.parent)
	}

	return nil

}

func AssignValue(name string, value interface{}, currentEnvironment *Environment) {
	if _, ok := currentEnvironment.values[name]; ok {
		currentEnvironment.values[name] = value
		return
	}

	if currentEnvironment.parent != nil {
		AssignValue(name, value, currentEnvironment.parent)
		return
	}

	panic("Undefined variable.")
}
