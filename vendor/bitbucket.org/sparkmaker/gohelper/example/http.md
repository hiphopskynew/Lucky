## HTTP library ##

### Example for using HTTP GET method ###

```go
package main

import (
	"time"
	"bitbucket.org/sparkmaker/gohelper/http"
)

func main() {
	response, err := http.Caller{
		URL:        "https://www.scale360solutions.com",
		RetryCount: 5,
		Headers: map[string]string{
			"Authentication": "xyz",
		},
	}.GET()
	...
}
```

### Example for using HTTP POST method ###

```go
package main

import (
	"time"
	"bitbucket.org/sparkmaker/gohelper/http"
)

func main() {
	response, err := http.Caller{
		URL:        "https://www.scale360solutions.com",
		RetryCount: 5,
		Headers: map[string]string{
			"Authentication": "xyz",
		},
		Body: "", // using interface type
	}.POST()
	...
}
```

### Example for using HTTP PUT method ###

```go
package main

import (
	"time"
	"bitbucket.org/sparkmaker/gohelper/http"
)

func main() {
	response, err := http.Caller{
		URL:        "https://www.scale360solutions.com",
		RetryCount: 5,
		Headers: map[string]string{
			"Authentication": "xyz",
		},
		Body: "", // using interface type
	}.PUT()
	...
}
```

### Example for using HTTP PATCH method ###

```go
package main

import (
	"time"
	"bitbucket.org/sparkmaker/gohelper/http"
)

func main() {
	response, err := http.Caller{
		URL:        "https://www.scale360solutions.com",
		RetryCount: 5,
		Headers: map[string]string{
			"Authentication": "xyz",
		},
		Body: "", // using interface type
	}.PATCH()
	...
}
```

### Example for using HTTP DELETE method ###

```go
package main

import (
	"time"
	"bitbucket.org/sparkmaker/gohelper/http"
)

func main() {
	response, err := http.Caller{
		URL:        "https://www.scale360solutions.com",
		RetryCount: 5,
		Headers: map[string]string{
			"Authentication": "xyz",
		},
		Body: "", // using interface type
	}.DELETE()
	...
}
```

**Remark: Fixed timeout 20 seconds can will configuration in the future** 