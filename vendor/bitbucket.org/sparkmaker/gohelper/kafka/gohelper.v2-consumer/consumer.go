package consumer

import (
	"encoding/json"
	"os"
	"os/signal"
	"strings"

	"bitbucket.org/sparkmaker/gohelper/logger/stdout"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

var memory map[string][]segment = make(map[string][]segment)

type segment struct {
	ID      string `json:"s_id"`
	Order   int    `json:"s_order"`
	Overall int    `json:"s_overall"`
	Data    string `json:"s_data"`
}

func (s segment) IsValid() bool {
	return s.ID != "" && s.Order != 0 && s.Overall != 0 && s.Data != ""
}

type HandleSuccess func(Message)
type HandleError func(error)

type Property struct {
	BrokerList            []string
	Topics                []string
	GroupID               string
	Config                *cluster.Config
	DisableAutoMarkOffset bool // default: false (auto mark offset when received message)
	DisableAutoRestart    bool // default: false (auto restart always)
	DisableLog            bool // default: false (enabled)
}

type Message struct {
	consumer        *cluster.Consumer
	consumerMessage *sarama.ConsumerMessage
	Key, Value      []byte
	Partition       int32
	Offset          int64
	Topic           string
}

func (msg Message) MarkOffset() {
	msg.consumer.MarkOffset(msg.consumerMessage, string(``))
}

func (msg Message) isSegment() bool {
	s := segment{}
	json.Unmarshal(msg.Value, &s)
	return s.IsValid()
}

type KafkaConsumer interface {
	Handle(HandleSuccess, HandleError)
	Terminate()
	IsConsuming() bool
}

type kafkaConsumerImp struct {
	consumer *cluster.Consumer
	status   bool
	property Property
}

func stringSerialized(parts []segment) string {
	var data = make([]string, parts[0].Overall, parts[0].Overall)
	for _, s := range parts {
		data[s.Order-1] = s.Data
	}
	return strings.Join(data, string(``))
}

func chunk(hs HandleSuccess) func(msg Message) {
	return func(msg Message) {
		// store segment of message to memory
		part := segment{}
		json.Unmarshal(msg.Value, &part)

		// serialization string message if full part and clear data in the memory, and send to process success function
		memory[part.ID] = append(memory[part.ID], part)
		if len(memory[part.ID]) == part.Overall && part.Order > 0 && part.Overall > 0 {
			msg.Value = []byte(stringSerialized(memory[part.ID]))
			hs(msg)
			delete(memory, part.ID)
		}
	}
}

func New(prop Property) KafkaConsumer {
	return &kafkaConsumerImp{property: prop}
}

func (kci *kafkaConsumerImp) IsConsuming() bool {
	return kci.status
}

func (kci *kafkaConsumerImp) Terminate() {
	kci.consumer.Close()
}

// Logging
func (kci *kafkaConsumerImp) info(msgs ...interface{}) {
	if kci.property.DisableLog {
		return
	}
	stdout.Info(msgs)
}

// Logging
func (kci *kafkaConsumerImp) error(msgs ...interface{}) {
	if kci.property.DisableLog {
		return
	}
	stdout.Error(msgs)
}

func (kci *kafkaConsumerImp) Handle(hs HandleSuccess, he HandleError) {
	// trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	doneCh := make(chan struct{})

	// consume messages, watch signals async
	go func(ch chan struct{}) {
		started := func() {
			kci.status = true
			kci.info("the kafka consumer started")
		}
		closed := func() {
			kci.status = false
			kci.info("the kafka consumer stopped")
		}

		config := kci.property.Config
		// setting default config
		if config == nil {
			config = cluster.NewConfig()
			config.Consumer.Return.Errors = true
		}

		// init consumer
		consumer, err := cluster.NewConsumer(kci.property.BrokerList, kci.property.GroupID, kci.property.Topics, config)
		if err != nil {
			kci.error(err.Error())
			panic(err)
		}

		kci.consumer = consumer
		defer func() {
			consumer.Close()
			if !kci.property.DisableAutoRestart {
				kci.info("the kafka consumer restarting")
				New(kci.property).Handle(hs, he)
			}
		}()

		go func() {
			started()
			defer func() {
				if err := recover(); err != nil {
					kci.error("the kafka consumer stopping because:", err)
				}
				closed()
			}()
			for {
				select {
				case err := <-consumer.Errors():
					he(err)
				case msg, ok := <-consumer.Messages():
					if !ok {
						continue
					}
					cm := Message{
						consumer:        consumer,
						consumerMessage: msg,
						Key:             msg.Key,
						Value:           msg.Value,
						Partition:       msg.Partition,
						Offset:          msg.Offset,
						Topic:           msg.Topic,
					}
					if !kci.property.DisableAutoMarkOffset {
						cm.MarkOffset()
					}
					if cm.isSegment() {
						chunk(hs)(cm)
						continue
					}
					hs(cm)
				case <-signals:
					ch <- struct{}{}
					return
				}
			}
		}()
		<-ch
	}(doneCh)
}
