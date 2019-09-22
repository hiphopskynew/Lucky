package consumer

import (
	"encoding/json"
	"os"
	"os/signal"
	"strings"

	cluster "github.com/bsm/sarama-cluster"
)

type KafkaConsumer struct {
	BrokerList []string
	Topics     []string
	Partition  int32 // unnecessary to set (the process will ignore this setting)
	GroupID    string
}

type Message struct {
	Key, Value []byte
	Topic      string
}

type HandleSuccess func(Message)
type HandleError func(string)

func (kc KafkaConsumer) Handle(hs HandleSuccess, he HandleError) {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true

	// init consumer
	consumer, err := cluster.NewConsumer(kc.BrokerList, kc.GroupID, kc.Topics, config)

	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	// trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	doneCh := make(chan struct{})

	// consume messages, watch signals
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				he(err.Error())
			case msg, ok := <-consumer.Messages():
				if ok {
					consumer.MarkOffset(msg, "") // mark message as processed
					hs(Message{
						Topic: msg.Topic,
						Key:   msg.Key,
						Value: msg.Value,
					})
				}
			case <-signals:
				doneCh <- struct{}{}
				return
			}
		}
	}()
	<-doneCh
}

type Segment struct {
	ID      string `json:"s_id"`
	Order   int    `json:"s_order"`
	Overall int    `json:"s_overall"`
	Data    string `json:"s_data"`
}

var memory map[string][]Segment = make(map[string][]Segment)

func stringSerialized(segments []Segment) string {
	var data = make([]string, segments[0].Overall, segments[0].Overall)
	for _, s := range segments {
		data[s.Order-1] = s.Data
	}
	return strings.Join(data, "")
}

func chunk(hs HandleSuccess) func(msg Message) {
	return func(msg Message) {
		// store segment of message to memory
		segment := Segment{}
		json.Unmarshal(msg.Value, &segment)

		// serialization string message if full segment and clear data in the memory, and send to process success function
		memory[segment.ID] = append(memory[segment.ID], segment)
		if len(memory[segment.ID]) == segment.Overall && segment.Order > 0 && segment.Overall > 0 {
			msg.Value = []byte(stringSerialized(memory[segment.ID]))
			hs(msg)
			delete(memory, segment.ID)
		}
	}
}

func (kc KafkaConsumer) HandleSegment(hs HandleSuccess, he HandleError) {
	kc.Handle(chunk(hs), he)
}
