package rule

type Failure struct {
	Key      string   `json:"key"`
	Messages []string `json:"messages"`
}
type Rule func(key string, value interface{}) *Failure

func error(key, message string) *Failure {
	return &Failure{Key: key, Messages: []string{message}}
}

func isValidRules(key string, value interface{}) func(...Rule) bool {
	return func(rules ...Rule) bool {
		for _, rule := range rules {
			if rule(key, value) != nil {
				return false
			}
		}
		return true
	}
}

func Required() Rule {
	return func(key string, value interface{}) *Failure {
		if value != nil {
			return nil
		}
		return error(key, "is required")
	}
}
