## Kafka Consumer Library ##

### Example for Kafka consumer (Full message) ###

```go
package main

import (
	"bitbucket.org/sparkmaker/gohelper/kafka/consumer"
)

func main() {
	consumer := consumer.KafkaConsumer{
		Topics:      []string{"topic"},
		BrokerList: []string{"localhost:9092"},
		GroupID: "test-group",
  }

  consumer.Handle(successFn, errorFn)
}

func successFn(msg consumer.Message) {
	...
}

func errorFn(str string) {
	...
}
```

### Example for Kafka consumer (Segment message) ###

```go
package main

import (
	"bitbucket.org/sparkmaker/gohelper/kafka/consumer"
)

func main() {
	consumer := consumer.KafkaConsumer{
		Topics:      []string{"topic"},
		BrokerList: []string{"localhost:9092"},
		GroupID: "test-group",
  }

  consumer.HandleSegment(successFn, errorFn)
}

func successFn(msg consumer.Message) {
	...
}

func errorFn(str string) {
	...
}
```

### Example for Kafka consumer version 2 (**Deprecated**) ###
This feature supported handle standalone message & segment messages in one handle

```go
package main

import (
	consumer "bitbucket.org/sparkmaker/gohelper/kafka/gohelper.v2-consumer"
	"bitbucket.org/sparkmaker/gohelper/logger/stdout"
)

func main() {
	property = consumer.Property{
		Topics:     []string{"gohelper.v2"},
		BrokerList: []string{"gohelper:9092"},
		GroupID:    "gohelper",
	}
	kcm := consumer.New(property)
	kcm.Handle(successFn, errorFn)
}

func successFn(msg consumer.Message) {
	...
}

func errorFn(err error) {
	...
}
```

### Example for Kafka consumer version 3 (**Recommended**) ###
This feature supported handle standalone message & segment messages in one handle

```go
package main

import (
	consumer "bitbucket.org/sparkmaker/gohelper/kafka/gohelper.v3-consumer"
	"bitbucket.org/sparkmaker/gohelper/logger/stdout"
)

func main() {
	property = consumer.Property{
		Topics:     []string{"gohelper.v3"},
		BrokerList: []string{"gohelper:9092"},
		GroupID:    "gohelper",
	}
	kcm := consumer.New(property)
	kcm.Handle(successFn, errorFn)
}

func successFn(msg consumer.Message) {
	...
}

func errorFn(err error) {
	...
}
```