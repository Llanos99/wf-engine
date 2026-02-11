package domain

type ActionSpec struct {
	ActionName string                 `yaml:"action_name" json:"action_name"`
	Input      map[string]interface{} `yaml:"input" json:"input"`
	Next       string                 `yaml:"next" json:"next"`
}

type IfSpec struct {
	Expression Expression `yaml:"expression" json:"expression"`
	TrueNext   string     `yaml:"true_next" json:"true_next"`
	FalseNext  string     `yaml:"false_next" json:"false_next"`
}

type WaitSpec struct {
	Duration string `yaml:"duration" json:"duration"`
	Next     string `yaml:"next" json:"next"`
}
