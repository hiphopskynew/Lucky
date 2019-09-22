package rule

import "fmt"

func IsList() Rule {
	return func(key string, value interface{}) *Failure {
		switch value.(type) {
		case []interface{}, nil:
			return nil
		default:
			return error(key, "is invalid type of list")
		}
	}
}

func NonEmptyList() Rule {
	return func(key string, value interface{}) *Failure {
		if !isValidRules(key, value)(Required(), IsList()) {
			return nil
		}
		if len(value.([]interface{})) == 0 {
			return error(key, "is empty")
		}
		return nil
	}
}

func MaxElement(n int) Rule {
	return func(key string, value interface{}) *Failure {
		if !isValidRules(key, value)(Required(), IsList()) {
			return nil
		}
		if len(value.([]interface{})) > n {
			return error(key, fmt.Sprintf("is more than %v elements", n))
		}
		return nil
	}
}

func MinElement(n int) Rule {
	return func(key string, value interface{}) *Failure {
		if !isValidRules(key, value)(Required(), IsList()) {
			return nil
		}
		if len(value.([]interface{})) < n {
			return error(key, fmt.Sprintf("is less than %v elements", n))
		}
		return nil
	}
}
