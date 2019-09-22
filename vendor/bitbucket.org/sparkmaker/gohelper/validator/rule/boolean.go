package rule

func IsBoolean() Rule {
	return func(key string, value interface{}) *Failure {
		switch value.(type) {
		case bool, nil:
			return nil
		default:
			return error(key, "is invalid type of boolean")
		}
	}
}
