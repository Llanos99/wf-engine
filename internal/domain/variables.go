package domain

import "fmt"

type VariableType int

type Variables map[string]interface{}

func NewVariables() Variables {
	return make(Variables)
}

func (v Variables) SetBool(key string, value bool) {
	v[key] = value
}

func (v Variables) SetString(key string, value string) {
	v[key] = value
}

func (v Variables) SetFloat(key string, value float64) {
	v[key] = value
}

func (v Variables) SetArray(key string, value []interface{}) {
	v[key] = value
}

func (v Variables) GetBool(key string) (bool, error) {
	val, ok := v[key]
	if !ok {
		return false, fmt.Errorf("variable %s not found", key)
	}
	b, ok := val.(bool)
	if !ok {
		return false, fmt.Errorf("variable %s is not a bool", key)
	}
	return b, nil
}

func (v Variables) GetString(key string) (string, error) {
	val, ok := v[key]
	if !ok {
		return "", fmt.Errorf("variable %s not found", key)
	}
	str, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("variable %s is not a string", key)
	}
	return str, nil
}

func (v Variables) GetFloat(key string) (float64, error) {
	val, ok := v[key]
	if !ok {
		return 0, fmt.Errorf("variable %s not found", key)
	}
	switch n := val.(type) {
	case float64:
		return n, nil
	case int:
		return float64(n), nil
	default:
		return 0, fmt.Errorf("variable %s is not a number", key)
	}
}

func (v Variables) GetArray(key string) ([]interface{}, error) {
	val, ok := v[key]
	if !ok {
		return nil, fmt.Errorf("variable %s not found", key)
	}
	arr, ok := val.([]interface{})
	if !ok {
		return nil, fmt.Errorf("variable %s is not an array", key)
	}
	return arr, nil
}

func (v Variables) Get(key string) (interface{}, bool) {
	val, ok := v[key]
	return val, ok
}
