## Service Level Agreement library ##

### SLA Configuration ###

| ENV                   | type   | default | description           |
| --------------------- | ------ | ------- | --------------------- |
| SLA_KAFKA_TOPIC       | string |         | SLA kafka topic       |
| SLA_KAFKA_MAX_RETRY   | int    | 0       | SLA kafka max retry   |
| SLA_KAFKA_KEY         | string |         | SLA kafka message key |
| SLA_KAFKA_BROKER_LIST | string |         | SLA kafka broker list |
| SLA_SERVICE_NAME      | string |         | Service name          |
| SLA_SERVICE_GROUP     | string |         | Service group         |
| SLA_VERSION           | string |         | SLA version           |
| SLA_INSTANCE_ID       | string |         | Docker instance ID    |
| SLA_SOURCE            | string |         | SLA source            |

### Example ###

```go
package middleware

import "bitbucket.org/sparkmaker/gohelper/logger/sla"

// In the middleware
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sla := sla.New(w, r)                    // New SLA
		next.ServeHTTP(sla.ResponseWriter(), r) // Forward response writer (using `sla.ResponseWriter()` instead of `w`)
		sla.Finished()                          // End of process
	})
}
```