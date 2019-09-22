package rule

import (
	"encoding/json"
	"strings"
)

func IsObject() Rule {
	return func(key string, value interface{}) *Failure {
		if !isValidRules(key, value)(Required()) {
			return nil
		}
		b, _ := json.Marshal(value)
		if strings.HasPrefix(string(b), "{") && strings.HasSuffix(string(b), "}") {
			return nil
		}
		return error(key, "is invalid type of object")
	}
}

func NonEmptyObject() Rule {
	return func(key string, value interface{}) *Failure {
		if !isValidRules(key, value)(Required(), IsObject()) {
			return nil
		}
		b, _ := json.Marshal(value)
		if string(b) == "{}" {
			return error(key, "is empty")
		}
		return nil
	}
}
