package models

import "fmt"

type VariableType int

const (
	BOOL VariableType = iota
	INT
	STRING
	LIST
)

type Variables struct {
	Data map[string]Variable
}

type Variable struct {
	Type  VariableType
	Value any
}

func (v *Variables) SetInt(key string, val int) {
	v.Data[key] = Variable{Type: INT, Value: val}
}

func (v *Variables) SetBool(key string, val bool) {
	v.Data[key] = Variable{Type: BOOL, Value: val}
}

func (v *Variables) SetString(key string, val string) {
	v.Data[key] = Variable{Type: STRING, Value: val}
}

func (v *Variables) GetInt(key string) (int, error) {
	variable, ok := v.Data[key]
	if !ok {
		return 0, fmt.Errorf("variable %s not found", key)
	}
	if variable.Type != INT {
		return 0, fmt.Errorf("variable %s is not int", key)
	}
	return variable.Value.(int), nil
}

func (v *Variables) GetBool(key string) (bool, error) {
	variable, ok := v.Data[key]
	if !ok {
		return false, fmt.Errorf("variable %s not found", key)
	}
	if variable.Type != BOOL {
		return false, fmt.Errorf("variable %s is not bool", key)
	}
	return variable.Value.(bool), nil
}

func (v *Variables) GetString(key string) (string, error) {
	variable, ok := v.Data[key]
	if !ok {
		return "", fmt.Errorf("variable %s not found", key)
	}
	if variable.Type != STRING {
		return "", fmt.Errorf("variable %s is not string", key)
	}
	return variable.Value.(string), nil
}
