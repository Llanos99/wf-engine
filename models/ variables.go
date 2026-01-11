package models

type VariableType int

const (
	BOOL VariableType = iota
	NUM
	STRING
	LIST
)

type Variables struct {
	data map[string]Variable
}

type Variable struct {
	Type  VariableType
	Value any
}
