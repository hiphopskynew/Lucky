## JWT Library ##

### Example for encode raw data to JWT token string ###

- **Encode without expire**

```go
package main

import (
	"encoding/json"
	"bitbucket.org/sparkmaker/gohelper/jwt"
)

type User struct {
	Firstname string
	Lastname  string
	Age       int
}

func main() {
	codec := jwt.Config("secret")
	tokenstring := codec.Encode(User{
		Firstname: "John",
		Lastname:  "Smith",
		Age:       30,
	})
	...
}

```

- **Encode with expire time**

```go
package main

import (
	"encoding/json"
	"bitbucket.org/sparkmaker/gohelper/jwt"
)

type User struct {
	Firstname string
	Lastname  string
	Age       int
}

func main() {
	codec := jwt.ConfigWithExpire("secret", time.Second)
	tokenstring := codec.Encode(User{
		Firstname: "John",
		Lastname:  "Smith",
		Age:       30,
	})
	...
}
```

### Example for decode JWT token ###

- **Encode JWT token and parse to struct**

```go
package main

import (
	"encoding/json"
	"bitbucket.org/sparkmaker/gohelper/jwt"
)

type User struct {
	Firstname string
	Lastname  string
	Age       int
}

func main() {
	codec := jwt.Config("secret")
  payload, err := codec.Decode(tokenstring) // A 'tokenstring' variable from before example

  ... // Handle error

	user := User{}
  json.Unmarshal(payload, &user) // A 'user' variable will be parse payload of JWT token to struct
	...
}

```