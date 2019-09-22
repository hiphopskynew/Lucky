package rule

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func floatToStr(n float64) string {
	split := strings.Split(fmt.Sprintf("%v", n), ".")
	regex, _ := regexp.Compile(`0+$`)
	digit := split[0]
	digitPoint := ""
	if len(split) > 1 {
		digitPoint = regex.ReplaceAllString(split[1], "")
	}
	result := digit
	if len(digitPoint) > 0 {
		if strings.Contains(digitPoint, "e+") {
			splitE := strings.Split(digitPoint, "e+")
			powerNum, _ := strconv.Atoi(splitE[1])
			digit = strings.Join([]string{digit, splitE[0][:powerNum]}, "")
			digitPoint = splitE[0][powerNum:]
		}
		result = digit
		if len(digitPoint) > 0 {
			result = strings.Join([]string{digit, digitPoint}, ".")
		}
	}
	return result
}

func IsNumeric() Rule {
	return func(key string, value interface{}) *Failure {
		switch value.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64, nil:
			return nil
		default:
			return error(key, "is invalid type of number")
		}
	}
}

func MinValue(n float64) Rule {
	return func(key string, value interface{}) *Failure {
		if !isValidRules(key, value)(Required(), IsNumeric()) {
			return nil
		}
		if value.(float64) < n {
			return error(key, fmt.Sprintf("value less than %v", floatToStr(n)))
		}
		return nil
	}
}

func EqualValue(n float64) Rule {
	return func(key string, value interface{}) *Failure {
		if !isValidRules(key, value)(Required(), IsNumeric()) {
			return nil
		}
		if value.(float64) != n {
			return error(key, fmt.Sprintf("value not equal to %v", floatToStr(n)))
		}
		return nil
	}
}

func MaxValue(n float64) Rule {
	return func(key string, value interface{}) *Failure {
		if !isValidRules(key, value)(Required(), IsNumeric()) {
			return nil
		}
		if value.(float64) > n {
			return error(key, fmt.Sprintf("value greater than %v", floatToStr(n)))
		}
		return nil
	}
}

func MaxPrecision(n int) Rule {
	return func(key string, value interface{}) *Failure {
		if !isValidRules(key, value)(Required(), IsNumeric()) {
			return nil
		}
		split := strings.Split(floatToStr(value.(float64)), ".")
		if len(split) > 1 && len(split[1]) > n {
			return error(key, fmt.Sprintf("value greater than %v precision", n))
		}
		return nil
	}
}
