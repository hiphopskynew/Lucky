package validator

import (
	"encoding/json"
	"strings"

	"bitbucket.org/sparkmaker/gohelper/validator/flatten"

	"bitbucket.org/sparkmaker/gohelper/validator/rule"
)

type DependFunc = func(map[string]interface{}) bool

type conditional struct {
	keys []string
	fns  []DependFunc
}

type validator struct {
	jsonStr  string
	jsonMap  map[string]interface{}
	rules    map[string][]rule.Rule
	errors   []rule.Failure
	dependOn map[string]conditional
}

func New(js string) *validator {
	v := &validator{jsonStr: js, rules: make(map[string][]rule.Rule)}
	njs, err := flatten.FlattenString(js, "")
	if err != nil {
		v.errors = append(v.errors, rule.Failure{Key: "", Messages: []string{"malformed json"}})
	}
	jm := make(map[string]interface{})
	err = json.Unmarshal([]byte(njs), &jm)
	v.jsonMap = jm
	if err != nil {
		v.errors = append(v.errors, rule.Failure{Key: "", Messages: []string{"malformed json"}})
	}
	v.dependOn = map[string]conditional{}
	return v
}

func cvtMsgList(fs []rule.Failure) []rule.Failure {
	result := []rule.Failure{}
	for _, f := range fs {
		if strings.Contains(f.Key, "$") {
			msgs := []string{}
			for _, msg := range f.Messages {
				if strings.HasPrefix(msg, "is invalid type of") {
					msgs = append(msgs, "is invalid type of list")
					continue
				}
				msgs = append(msgs, msg)
			}
			f.Messages = msgs
		}
		result = append(result, f)
	}
	return result
}

func aggregate(fs []rule.Failure) []rule.Failure {
	result := []rule.Failure{}
	unique := func(l []string) []string {
		m := make(map[string]interface{})
		r := []string{}
		for _, v := range l {
			m[v] = nil
		}
		for k := range m {
			r = append(r, k)
		}
		return r
	}
	find := func(e rule.Failure) int {
		for i, v := range result {
			if v.Key == e.Key {
				return i
			}
		}
		return -1
	}
	for _, f := range fs {
		rIndex := find(f)
		if rIndex != -1 {
			result[rIndex].Messages = unique(append(result[rIndex].Messages, f.Messages...))
			continue
		}
		result = append(result, f)
	}
	return result
}

func checkDependOn(fs []rule.Failure, v *validator) []rule.Failure {
	dependOn := v.dependOn
	newFailures := []rule.Failure{}
	isError := func(key string) bool {
		for _, f := range fs {
			if f.Key == key {
				return true
			}
		}
		return false
	}
	for _, f := range fs {
		doks, ok := dependOn[f.Key]
		if !ok {
			newFailures = append(newFailures, f)
			continue
		}
		addTag := true
		for _, dok := range doks.keys {
			_, keyExist := v.jsonMap[dok]
			if isError(dok) || !keyExist {
				addTag = false
			}
		}
		if !addTag {
			continue
		}
		for _, fn := range doks.fns {
			if !fn(v.jsonMap) {
				addTag = false
				break
			}
		}
		if addTag {
			newFailures = append(newFailures, f)
		}
	}
	return newFailures
}

func (v *validator) AddRule(key string, rules ...rule.Rule) {
	v.rules[key] = append(v.rules[key], rules...)
}

func (v *validator) DependOn(key string, dependOn []string, fns ...DependFunc) {
	v.dependOn[key] = conditional{keys: dependOn, fns: fns}
}

func (v *validator) Validate() []rule.Failure {
	if len(v.errors) > 0 {
		return v.errors
	}

	validate := func(key string, value interface{}, rs []rule.Rule) {
		failure := rule.Failure{Key: key, Messages: []string{}}
		for _, r := range rs {
			if err := r(key, value); err != nil {
				failure.Messages = append(failure.Messages, err.Messages...)
			}
		}
		if len(failure.Messages) > 0 {
			v.errors = append(v.errors, failure)
		}
	}
	for key, rs := range v.rules {
		_, keyExist := v.jsonMap[key]
		if strings.Contains(key, "$") && keyExist {
			for _, sliceV := range v.jsonMap[key].([]interface{}) {
				validate(key, sliceV, rs)
			}
		} else {
			validate(key, v.jsonMap[key], rs)
		}
	}
	return checkDependOn(cvtMsgList(aggregate(v.errors)), v)
}
