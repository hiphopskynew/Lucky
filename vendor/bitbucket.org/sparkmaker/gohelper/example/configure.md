## Configure library ##

The configuration will be auto replacing by environment variables first this is an example

In the configuration file `.json`

```json
{
	"a": {
		"b": {
			"c": "value"
		}
	}
}
```

The environment variable upper case only for replacing the configuration file and tear down JSON level by underscore e.g. `A_B_C` will be replacing value if you getting `a.b.c` example in the configuration file above.

### How to set run mode in the environment variable ###
This code for setting run mode

```sh
export RUN_MODE="production"
```

### Example for getting a string in the configuration file ###

```go
package main

import "bitbucket.org/sparkmaker/gohelper/configure"

func main() {
	config := configure.New() // Default configuration file path is `/conf/application.json` and run mode `default`
	config.GetString("a.b.c")
}
```

### Example for getting a integer in the configuration file ###

```go
package main

import "bitbucket.org/sparkmaker/gohelper/configure"

func main() {
	config := configure.New() // Default configuration file path is `/conf/application.json` and run mode `default`
	config.GetInt("a.b.c")
}
```

### Example for getting a boolean in the configuration file ###

```go
package main

import "bitbucket.org/sparkmaker/gohelper/configure"

func main() {
	config := configure.New() // Default configuration file path is `/conf/application.json` and run mode `default`
	config.GetBool("a.b.c")
}
```

### Example for getting a list of string in the configuration file ###

```go
package main

import "bitbucket.org/sparkmaker/gohelper/configure"

func main() {
	config := configure.New() // Default configuration file path is `/conf/application.json` and run mode `default`
	config.GetStrings("a.b.c")
}
```

### Example for adding a configuration into run mode ###

```go
package main

import "bitbucket.org/sparkmaker/gohelper/configure"

func main() {
	config := configure.New()
	config.Add("production", "/conf/production.json")
	config.Add("develop", "/conf/develop.json")
	config.Add("test", "/conf/test.json")

	// strict mode will need to set all variables via environment
	config.StrictMode([]string{"production"}, []string{})
	...
}
```

### Example for delete configuration from run mode ###

```go
package main

import "bitbucket.org/sparkmaker/gohelper/configure"

func main() {
	config := configure.New()
	config.Del("test")
	...
}
```