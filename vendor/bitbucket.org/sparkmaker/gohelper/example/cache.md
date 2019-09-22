## Cache library ##

### Example for set without expire into memory cache ###

```go
package main

import "bitbucket.org/sparkmaker/gohelper/cache"

func main() {
	err := cache.Set("k", map[string]interface{}{})
	...
}
```

### Example for set with expire into memory cache ###

```go
package main

import (
	"time"
	"bitbucket.org/sparkmaker/gohelper/cache"
)

func main() {
	err := cache.SetWithExpire("k", map[string]interface{}{}, time.Hour)
	...
}
```

### Example for get memory cache and can custom handle ###

```go
package main

import "bitbucket.org/sparkmaker/gohelper/cache"

func main() {
	value, err := cache.Get("key")
	...
}
```

### Example for get memory cache and recover when memory cache exception occur ###

```go
package main

import "bitbucket.org/sparkmaker/gohelper/cache"

func main() {
	value := cache.GetWithoutErr("key")
	...
}
```