## XHeader integration with HTTP Caller library ##

### Example for parsing request header to XHeader ###

```go
package main

import (
	"fmt"

	"bitbucket.org/sparkmaker/gohelper/http"
)

func main() {
	reqHeader := map[string]string{
		"Content-Type": "application/json",
		"x-one":      	"One direction",
		"x-two":      	"Two direction",
	}
	xh := http.NewXHeader()
	xh.Append(reqHeader)
	...
}
```

### Example for additional XHeader ###

```go
package main

import (
	"fmt"

	"bitbucket.org/sparkmaker/gohelper/http"
)

func main() {
	...
	xh.Add("x-three", "Three direction")
	...
}
```

### Example for remove key in XHeader ###

```go
package main

import (
	"fmt"

	"bitbucket.org/sparkmaker/gohelper/http"
)

func main() {
	...
	xh.Del("x-one")
	...
}
```

### Example for parsing XHeader to map[string]string type ###

```go
package main

import (
	"fmt"

	"bitbucket.org/sparkmaker/gohelper/http"
)

func main() {
	...
	m := xh.ToMap()
	...
}
```

### Example for custom prefix of XHeader (Default is 'x-') ###

```go
package main

import (
	"fmt"

	"bitbucket.org/sparkmaker/gohelper/http"
)

func main() {
	...
	xh.SetPrefix("z-")
	...
}
```
**Remark** If you want to custom prefix. You should setting after new XHeader.