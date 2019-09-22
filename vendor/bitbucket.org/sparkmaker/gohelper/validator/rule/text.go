package rule

import (
	"fmt"
	"regexp"
	"strings"
)

func IsString() Rule {
	return func(key string, value interface{}) *Failure {
		switch value.(type) {
		case string, nil:
			return nil
		default:
			return error(key, "is invalid type of string")
		}
	}
}

func NonEmpty() Rule {
	return func(key string, value interface{}) *Failure {
		if !isValidRules(key, value)(Required(), IsString()) {
			return nil
		}
		if len([]rune(strings.TrimSpace(value.(string)))) == 0 {
			return error(key, "is empty")
		}
		return nil
	}
}

func MinLength(n int) Rule {
	return func(key string, value interface{}) *Failure {
		if !isValidRules(key, value)(Required(), IsString()) {
			return nil
		}

		if len([]rune(strings.TrimSpace(value.(string)))) < n {
			return error(key, fmt.Sprintf("is less than %v characters", n))
		}
		return nil
	}
}

func MaxLength(n int) Rule {
	return func(key string, value interface{}) *Failure {
		if !isValidRules(key, value)(Required(), IsString()) {
			return nil
		}
		if len([]rune(strings.TrimSpace(value.(string)))) > n {
			return error(key, fmt.Sprintf("is more than %v characters", n))
		}
		return nil
	}
}

func EqualLength(n int) Rule {
	return func(key string, value interface{}) *Failure {
		if !isValidRules(key, value)(Required(), IsString()) {
			return nil
		}
		if len([]rune(strings.TrimSpace(value.(string)))) != n {
			return error(key, fmt.Sprintf("is not equal to %v characters", n))
		}
		return nil
	}
}

func Format(expr string) Rule {
	return func(key string, value interface{}) *Failure {
		if !isValidRules(key, value)(Required(), IsString()) {
			return nil
		}
		r, err := regexp.Compile(expr)
		if err != nil {
			return error(key, err.Error())
		}
		if !r.Match([]byte(strings.TrimSpace(value.(string)))) {
			return error(key, "is invalid value format")
		}
		return nil
	}
}

func ExistIn(ls []string) Rule {
	return func(key string, value interface{}) *Failure {
		if !isValidRules(key, value)(Required(), IsString()) {
			return nil
		}
		trimedValue := strings.TrimSpace(value.(string))
		for _, v := range ls {
			if trimedValue == v {
				return nil
			}
		}
		return error(key, fmt.Sprintf("does not contain [%s]", strings.Join(ls, ", ")))
	}
}
