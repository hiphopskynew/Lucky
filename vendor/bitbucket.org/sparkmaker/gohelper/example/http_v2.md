# HTTP Caller V2 library #
HTTP Caller v2. There are a number of significant updates in this version that we hope you will like

1. Fixed timeout problem (http caller timeout)
2. Fixed http client cache problem
3. Adding retry setting you can see more detail in table below
4. Adding type retry function that can make you develop your own function for retry condition
5. Adding Printlog setting for print the log of library


## AttemRetry Settings ##
| Field Name  | Required |            Type             |      Default Value      |                                 Description                                  |
| :---------: | :------: | :-------------------------: | :---------------------: | :--------------------------------------------------------------------------: |
| attempRetry |  false   |             int             |      default is 1       |                     attempRetry for counting retry total                     |
|  MaxRetry   |  false   |             int             | default is 0 (no retry) |                         maximum of retry http caller                         |
|   Timeout   |  false   |        time.Duration        |  default is 30 seconds  |                     maximum time of http client timeout                      |
|  TimeDelay  |  false   |        time.Duration        |  default is 5 seconds   |    time delay between request (timeout + timedelay = next request coming)    |
|   RetryFn   |  false   | function(AttemRetry) bool |       no default        |                         function for retry condition                         |
|  RetryFns   |  false   |          []RetryFn          |       no default        |    list of retry function using for run mutiple validate retry functions     |
|  Response   |  false   |       *http.Response        |     default is nil      | add http response to this field for using when we want to add retry function |

### Example for using HTTP GET method ###

```go
package main

import (
	"time"
	"bitbucket.org/sparkmaker/gohelper/gohelper.v2-http"
)

func main() {
	response, err := http.New(http.Property{
		URL:        "https://www.scale360solutions.com",
		Headers: map[string]string{
			"Authentication": "xyz",
        },
        PrintLog: true,
        AttempRetry: &http.AttempRetry{
			TimeDelay: 5 * time.Second,
			Timeout:   3 * time.Second,
			MaxRetry:  5,
			RetryFns:  []http.RetryFn{http.RetryWithStatuses(400,500,503)},
        },
	}).GET()
	...
}
```

### Example for using HTTP POST method ###

```go
package main

import (
	"time"
	"bitbucket.org/sparkmaker/gohelper/gohelper.v2-http"
)

func main() {
	response, err := http.New(http.Property{
		URL:        "https://www.scale360solutions.com",
		Headers: map[string]string{
			"Authentication": "xyz",
        },
        PrintLog: true,
        Body: "", // using interface type
        AttempRetry: &http.AttempRetry{
			TimeDelay: 5 * time.Second,
			Timeout:   3 * time.Second,
			MaxRetry:  5,
			RetryFns:  []http.RetryFn{http.RetryWithStatuses(400,500,503)},
		},
	}).POST()
	...
}
```

### Example for using HTTP PUT method ###

```go
package main

import (
	"time"
	"bitbucket.org/sparkmaker/gohelper/gohelper.v2-http"
)

func main() {
	response, err := http.New(http.Property{
		URL:        "https://www.scale360solutions.com",
		Headers: map[string]string{
			"Authentication": "xyz",
        },
        PrintLog: true,
        Body: "", // using interface type
        AttempRetry: &http.AttempRetry{
			TimeDelay: 5 * time.Second,
			Timeout:   3 * time.Second,
			MaxRetry:  5,
			RetryFns:  []http.RetryFn{http.RetryWithStatuses(400,500,503)},
		},
	}).PUT()
	...
}
```

### Example for using HTTP PATCH method ###

```go
package main

import (
	"time"
	"bitbucket.org/sparkmaker/gohelper/gohelper.v2-http"
)

func main() {
	response, err := http.New(http.Property{
		URL:        "https://www.scale360solutions.com",
		Headers: map[string]string{
			"Authentication": "xyz",
        },
        PrintLog: true,
        Body: "", // using interface type
        AttempRetry: &http.AttempRetry{
			TimeDelay: 5 * time.Second,
			Timeout:   3 * time.Second,
			MaxRetry:  5,
			RetryFns:  []http.RetryFn{http.RetryWithStatuses(400,500,503)},
		},
	}).PATCH()
	...
}
```

### Example for using HTTP DELETE method ###

```go
package main

import (
	"time"
	"bitbucket.org/sparkmaker/gohelper/gohelper.v2-http"
)

func main() {
	response, err := http.New(http.Property{
		URL:        "https://www.scale360solutions.com",
		Headers: map[string]string{
			"Authentication": "xyz",
        },
        PrintLog: true,
        Body: "", // using interface type
        AttempRetry: &http.AttempRetry{
			TimeDelay: 5 * time.Second,
			Timeout:   3 * time.Second,
			MaxRetry:  5,
			RetryFns:  []http.RetryFn{http.RetryWithStatuses(400,500,503)},
		},
	}).DELETE()
	...
}
```

