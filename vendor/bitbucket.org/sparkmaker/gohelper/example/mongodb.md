## MongoDB [using github.com/globalsign/mgo] ##

This just wrap the library `github.com/globalsign/mgo` and manage connection.
The connection will be auto close after new a connection (Default 10 seconds) if forgot manual closing.

### Feature ###
- Custom connection limit.
- Custom time to close after new a connection.
- Custom waiting MongoDB response timeout.

- **New Connection without custom configuration**

```go
package main

import (
	"bitbucket.org/sparkmaker/gohelper/mongodb"
	"github.com/globalsign/mgo/bson"
)

func main() {
	url := "mongodb://localhost:27017"
	mongo := mongodb.New(url, "database_name")
	mongo.Database.C("collection_name").Insert(bson.M{"key": "value"})
	mongo.Close()
}
```

- **New Connection with custom configuration**

```go
package main

import (
	"bitbucket.org/sparkmaker/gohelper/mongodb"
	"github.com/globalsign/mgo/bson"
)

func main() {
	url := "mongodb://localhost:27017"
	mongo := mongodb.NewWithOption(url, "database_name", mongodb.Option{ConnectionLimitPerHost: 20})
	mongo.Database.C("collection_name").Insert(bson.M{"key": "value"})
	mongo.Close()
}
```

- **New Connection with credentials**

```go
package main

import (
	"bitbucket.org/sparkmaker/gohelper/mongodb"
	"github.com/globalsign/mgo/bson"
)

func main() {
	url := "mongodb://localhost:27017"
	mongo := mongodb.NewWithOption(url, "database_name", mongodb.Option{Username: "username", Password: "password"})
	mongo.Database.C("collection_name").Insert(bson.M{"key": "value"})
	mongo.Close()
}
```