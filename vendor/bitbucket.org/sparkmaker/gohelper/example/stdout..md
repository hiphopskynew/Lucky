## Standard output library ##

### Example for log INFO level ###

```go
package main

import "bitbucket.org/sparkmaker/gohelper/logger/stdout"

func main() {
	stdout.Info("message to log") // Can send multiple params for logging and interface type.
}
```

### Example for log DEBUG level ###

```go
package main

import "bitbucket.org/sparkmaker/gohelper/logger/stdout"

func main() {
	stdout.Debug("message to log") // Can send multiple params for logging and interface type.
}
```

### Example for log ERROR level ###

```go
package main

import "bitbucket.org/sparkmaker/gohelper/logger/stdout"

func main() {
	stdout.Error("message to log") // Can send multiple params for logging and interface type.
}
```

### Example for log TRACE level ###

```go
package main

import "bitbucket.org/sparkmaker/gohelper/logger/stdout"

func main() {
	stdout.Trace("message to log") // Can send multiple params for logging and interface type.
}
```

### Example for log Warning level ###

```go
package main

import "bitbucket.org/sparkmaker/gohelper/logger/stdout"

func main() {
	stdout.Warning("message to log") // Can send multiple params for logging and interface type.
}
```

## **Remark: Now this feature incomplete 100%. But you can using it.** ##