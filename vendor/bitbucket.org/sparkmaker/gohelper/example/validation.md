## Validation ##

This library using for validating JSON string, that now supported multiple level of JSON schema separate by "."
e.g. request JSON data
```json
{
	"a": {
		"b": [
			{
				"x": "X",
				"y": 0,
				"z": 10.5
			},
			{
				"x": "x",
				"y": 1,
				"z": 11.5
			}
		]
	}
}
```
example for validate y in the array index 0 using key "a.b.0.y"
example for validate y in the array index 1 using key "a.b.1.y"
example for validate y in the array all object using key "a.b.$.y" (**$** meaning is all of array index)

### **All rules of types provided** ###

**All types**
- Required()
- AnyExistIn([]interface{})

**Text type**
- IsString()
- NonEmpty()
- MinLength(int)
- MaxLength(int)
- EqualLength(int)
- Format(string)

**Numeric type**
- IsNumeric()
- MinValue(float64)
- EqualValue(float64)
- MaxValue(float64)
- MaxPrecision(int)

**Boolean type**
- IsBoolean()

**List type**
- IsList()
- NonEmptyList()
- MaxElement(int)
- MinElement(int)

**Object type**
- IsObject()
- NonEmptyObject()

**Remark** If using rule wrong type will skip that rule.

## This is an example of using ##

```go
package main

import (
	"bitbucket.org/sparkmaker/gohelper/validator"
	"bitbucket.org/sparkmaker/gohelper/validator/rule"
)

func main() {
	jsonStr := `{
		"keyStr": "string",
		"keyInt": 10
	}`

	validator := validator.New(jsonStr)
	validator.AddRule("keyStr", rule.IsString(), rule.MaxLength(10), rule.Format("^s.*g$"))
	validator.AddRule("keyInt", rule.IsNumeric(), rule.MinValue(1))

	// 'keyInt' will be validated when 'xKey' and 'yKey' is correctly.
	validator.DependOn("keyInt", []string{"xKey", "yKey"})

	// 'keyInt' will be validated when 'xKey', 'yKey', and modify function is correctly
	validator.DependOn("keyInt", []string{"xKey", "yKey"}, func(data map[string]interface{}) bool {
		if value, ok := data["xKey"].(string); ok {
			return value == "string"
		}
		return false
	})

	errors := validator.Validate()
	if len(errors) > 0 {
		...
	}
}

```