package domain

import "fmt"

func (e *Expression) Evaluate(vars Variables) (bool, error) {
	leftVal, ok := vars[e.Left]
	if !ok {
		return false, fmt.Errorf("variable '%s' not found", e.Left)
	}
	switch e.Operator {
	case "==":
		return leftVal == e.Right, nil
	case "!=":
		return leftVal != e.Right, nil
	case ">":
		return compareNumbers(leftVal, e.Right, ">")
	case "<":
		return compareNumbers(leftVal, e.Right, "<")
	case ">=":
		return compareNumbers(leftVal, e.Right, ">=")
	case "<=":
		return compareNumbers(leftVal, e.Right, "<=")
	default:
		return false, fmt.Errorf("unsupported operator: %s", e.Operator)
	}
}

func compareNumbers(left interface{}, right interface{}, op string) (bool, error) {
	lf, lok := toFloat(left)
	rf, rok := toFloat(right)
	if !lok || !rok {
		return false, fmt.Errorf("invalid numeric comparison")
	}
	switch op {
	case ">":
		return lf > rf, nil
	case "<":
		return lf < rf, nil
	case ">=":
		return lf >= rf, nil
	case "<=":
		return lf <= rf, nil
	}
	return false, fmt.Errorf("invalid operator")
}

func toFloat(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case float32:
		return float64(val), true
	case float64:
		return val, true
	default:
		return 0, false
	}
}
