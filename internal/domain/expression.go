package domain

type Operator string

const (
	OpEquals           Operator = "=="
	OpNotEquals        Operator = "!="
	OpGreaterThan      Operator = ">"
	OpLessThan         Operator = "<"
	OpGreaterThanEqual Operator = ">="
	OpLessThanEqual    Operator = "<="
)

type Expression struct {
	Left     string      `yaml:"left" json:"left"`
	Operator Operator    `yaml:"operator" json:"operator"`
	Right    interface{} `yaml:"right" json:"right"`
}
