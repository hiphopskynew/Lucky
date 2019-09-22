package flatten

import (
	"encoding/json"
	"errors"
	"strings"
)

// Nested input must be a map or slice
var NotValidInputError = errors.New("Not a valid input: map or slice")

func FlattenString(nestedstr, prefix string) (string, error) {
	var nested map[string]interface{}
	err := json.Unmarshal([]byte(nestedstr), &nested)
	if err != nil {
		return "", err
	}

	flatmap, err := FlattenAllPossibleKeys(nested, prefix)
	if err != nil {
		return "", err
	}

	flatb, err := json.Marshal(&flatmap)
	if err != nil {
		return "", err
	}

	return string(flatb), nil
}

func FlattenAllPossibleKeys(nested map[string]interface{}, prefix string) (map[string]interface{}, error) {
	flatMap := make(map[string]interface{})

	//1. when adding new key to flatMap, dupKeyCheck of the key is false
	//2. when creating list for children value, dupKeyCheck of the key is true
	dupKeyCheck := make(map[string]bool)

	err := flatten(true, flatMap, nested, prefix, dupKeyCheck)
	if err != nil {
		return nil, err
	}

	return flatMap, nil
}

func flatten(top bool, flatMap map[string]interface{}, nested interface{}, prefix string, dupKeyCheck map[string]bool) error {
	assign := func(newKey string, v interface{}) error {
		if ov, ok := flatMap[newKey]; ok {
			switch ov.(type) {
			case []interface{}:
				if c, ok := dupKeyCheck[newKey]; ok && !c {
					switch v.(type) {
					case []interface{}:
						nv := []interface{}{}
						flatMap[newKey] = append(nv, ov, v)
						dupKeyCheck[newKey] = true
					default:
						flatMap[newKey] = append(ov.([]interface{}), v)
					}
					break
				}
				flatMap[newKey] = append(ov.([]interface{}), v)
			default:
				nv := []interface{}{}
				flatMap[newKey] = append(nv, ov, v)
				dupKeyCheck[newKey] = true
			}
		} else {
			if strings.Contains(newKey, "$") {
				nv := []interface{}{}
				flatMap[newKey] = append(nv, v)
				dupKeyCheck[newKey] = true
			} else {
				flatMap[newKey] = v
				dupKeyCheck[newKey] = false
			}
		}

		switch v.(type) {
		case map[string]interface{}, []interface{}:
			if err := flatten(false, flatMap, v, newKey, dupKeyCheck); err != nil {
				return err
			}
		}
		return nil
	}

	switch nested.(type) {
	case map[string]interface{}:
		for k, v := range nested.(map[string]interface{}) {
			newKey := enkey(top, prefix, k)
			assign(newKey, v)
		}
	case []interface{}:
		for _, v := range nested.([]interface{}) {
			newKey := enkey(top, prefix, "$")
			assign(newKey, v)
		}
	default:
		return NotValidInputError
	}

	return nil
}

func enkey(top bool, prefix, subkey string) string {
	key := prefix
	if top {
		key += subkey
	} else {
		key += "." + subkey
	}
	return key
}
