## Kafka Producer ##

### Example for send message to the Kafka server (for small message) ###

```go
package main

import (
	"bitbucket.org/sparkmaker/gohelper/kafka/producer"
)

func main() {
	producer := producer.KafkaProducer{
		Topic:      "topic",
		MaxRetry:   5,
		BrokerList: []string{"localhost:9092"},
	}

	producer.Send("key", "message") // can using interface type for message
}
```

### Example for send message to the Kafka server (for large message) ###

```go
package main

import (
	"bitbucket.org/sparkmaker/gohelper/kafka/producer"
)

func main() {
	producer := &producer.KafkaProducer{
		Topic:      "topic",
		MaxRetry:   5,
		BrokerList: []string{"localhost:9092"},
	}

	producer.SetMaxSizeMessage(10)	// max message size per chunk (if not set, default as maximum 950000 bytes per chunk)
	producer.SendSegment("key", `{"some":"large message"}`) // only type 'string' supported for message to be sent
}
```

