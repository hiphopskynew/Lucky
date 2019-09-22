package rule

import (
	"encoding/json"
	"fmt"
)

// AnyExistIn supports type
// - int, int8, int16, int32, int64, float32, float64
// - string
// - bool
func AnyExistIn(ls []interface{}) Rule {
	return func(key string, value interface{}) *Failure {
		if !isValidRules(key, value)(Required()) {
			return nil
		}
		for _, v := range ls {
			switch v.(type) {
			case int:
				if int(value.(float64)) == v.(int) {
					return nil
				}
			case int8:
				if int8(value.(float64)) == v.(int8) {
					return nil
				}
			case int16:
				if int16(value.(float64)) == v.(int16) {
					return nil
				}
			case int32:
				if int32(value.(float64)) == v.(int32) {
					return nil
				}
			case int64:
				if int64(value.(float64)) == v.(int64) {
					return nil
				}
			case float32:
				if float32(value.(float64)) == v.(float32) {
					return nil
				}
			case float64:
				if value.(float64) == v.(float64) {
					return nil
				}
			default:
				if value == v {
					return nil
				}
			}
		}
		bytes, _ := json.Marshal(ls)
		return error(key, fmt.Sprintf("does not contain %s", string(bytes)))
	}
}
