package http

import (
	"net/http"
	"time"
)

// RetryFn type for create your own function retry
type RetryFn func(AttempRetry) bool

type AttempRetry struct {
	attempRetry int
	MaxRetry    int
	Timeout     time.Duration
	TimeDelay   time.Duration
	RetryFns    []RetryFn
	Response    *http.Response
}

// RetryWithStatuses retry when http status is one of these statuses.
func RetryWithStatuses(statusCodes ...int) RetryFn {
	return func(attemp AttempRetry) bool {
		for _, statusCode := range statusCodes {
			if statusCode == attemp.Response.StatusCode {
				return true
			}
		}
		return false
	}
}
