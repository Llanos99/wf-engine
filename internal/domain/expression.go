package domain

type Operator string

const (
	OpEquals      Operator = "eq"
	OpNotEquals   Operator = "neq"
	OpGreaterThan Operator = "gt"
	OpLessThan    Operator = "lt"
)

type Expression struct {
	Left     string      `yaml:"left" json:"left"`
	Operator Operator    `yaml:"operator" json:"operator"`
	Right    interface{} `yaml:"right" json:"right"`
}
