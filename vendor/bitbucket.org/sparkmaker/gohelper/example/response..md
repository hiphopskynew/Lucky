## Standard Response Template ##

### Normal function to make response wrapper ###

#### Success response group ####

- **OK response**

```go
package main

import "bitbucket.org/sparkmaker/gohelper/response"

func main() {
	respWrapper := response.OK(true) // Can using interface type
	...
}
```

- **Created response**

```go
package main

import "bitbucket.org/sparkmaker/gohelper/response"

func main() {
	respWrapper := response.CREATED(true) // Can using interface type
	...
}
```

#### Error response group ####

- **Bad Request response**

```go
package main

import "bitbucket.org/sparkmaker/gohelper/response"

func main() {
	respWrapper := response.BAD_REQUEST([]string{"Bad Request"}) // Can using interface type
	...
}
```

- **Unauthorized response**

```go
package main

import "bitbucket.org/sparkmaker/gohelper/response"

func main() {
	respWrapper := response.UNAUTHORIZED([]string{"Unauthorized"}) // Can using interface type
	...
}
```

- **Not Found response**

```go
package main

import "bitbucket.org/sparkmaker/gohelper/response"

func main() {
	respWrapper := response.NOT_FOUND([]string{"Not Found"}) // Can using interface type
	...
}
```

- **No Content response**

```go
package main

import "bitbucket.org/sparkmaker/gohelper/response"

func main() {
	respWrapper := response.GONE([]string{"Gone"}) // Can using interface type
	...
}
```

- **Internal Server Error response**

```go
package main

import "bitbucket.org/sparkmaker/gohelper/response"

func main() {
	respWrapper := response.INTERNAL_SERVER_ERROR([]string{"Internal Server Error"}) // Can using interface type
	...
}
```

- **Service Unavailable response**

```go
package main

import "bitbucket.org/sparkmaker/gohelper/response"

func main() {
	respWrapper := response.SERVICE_UNAVAILABLE([]string{"Service Unavailable"}) // Can using interface type
	...
}
```

### Custom function to make response ###

- **Success response**

```go
package main

import "bitbucket.org/sparkmaker/gohelper/response"

func main() {
	respWrapper := response.DataResponse(true, 301) // Can using interface type
	...
}
```

- **Error response**

```go
package main

import "bitbucket.org/sparkmaker/gohelper/response"

func main() {
	respWrapper := response.ErrorResponse("Message", []string{"Error info 1", "Error info 2"}, 99, 502)
	...
}
```

- **Custom Error response** 
 
```go 
package main 
 
import "bitbucket.org/sparkmaker/gohelper/response" 
 
func main() { 
  respWrapper := response.CustomErrorResponse("Error info", 502) 
  ... 
} 
```

### Response Internal Code ###

|Response									|code    |
| ----------------------- | ------ |
|Bad Request							|10			 |
|Unauthorized							|20      |
|Not Found								|30      |
|No Content 							|30			 |
|Internal Server Error		|40			 |
|Service Unavailable			|50		   |