package consumer

import (
	"context"
	"encoding/json"
	"strings"

	"bitbucket.org/sparkmaker/gohelper/logger/stdout"

	"github.com/Shopify/sarama"
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
	BrokerList             []string
	Topics                 []string
	GroupID                string
	Config                 *sarama.Config
	DisableAutoMarkMessage bool // default: false (auto mark message when received message)
	DisableAutoRestart     bool // default: false (auto restart always)
	DisableLog             bool // default: false (enabled)
}

type Message struct {
	consumerGroupClaim   sarama.ConsumerGroupClaim
	consumerGroupSession sarama.ConsumerGroupSession
	consumerMessage      *sarama.ConsumerMessage
	Key, Value           []byte
	Partition            int32
	Offset               int64
	Topic                string
}

func (msg Message) MarkMessage() {
	msg.consumerGroupSession.MarkMessage(msg.consumerMessage, string(``))
}

func (msg Message) isSegment() bool {
	s := segment{}
	json.Unmarshal(msg.Value, &s)
	return s.IsValid()
}

type Consumer struct {
	kafkaConsumer KafkaConsumer
	handleSuccess HandleSuccess
	handleError   HandleError
	ready         chan bool
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		cm := Message{
			consumerGroupClaim:   claim,
			consumerGroupSession: session,
			consumerMessage:      msg,
			Key:                  msg.Key,
			Value:                msg.Value,
			Partition:            msg.Partition,
			Offset:               msg.Offset,
			Topic:                msg.Topic,
		}
		if !consumer.kafkaConsumer.Property().DisableAutoMarkMessage {
			cm.MarkMessage()
		}
		if cm.isSegment() {
			chunk(consumer.handleSuccess)(cm)
			continue
		}
		consumer.handleSuccess(cm)
	}
	return nil
}

type KafkaConsumer interface {
	Handle(HandleSuccess, HandleError)
	Terminate()
	IsConsuming() bool
	Property() Property
}

type kafkaConsumerImp struct {
	consumerGroup sarama.ConsumerGroup
	status        bool
	property      Property
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
	kci.consumerGroup.Close()
}

func (kci *kafkaConsumerImp) Property() Property {
	return kci.property
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
	// consume messages, watch signals async
	started := func() {
		kci.status = true
		kci.info("the kafka consumer started")
	}
	closed := func() {
		kci.status = false
		kci.info("the kafka consumer stopped")
	}

	version, _ := sarama.ParseKafkaVersion("2.1.1")
	config := kci.property.Config
	// setting default config
	if config == nil {
		config = sarama.NewConfig()
		config.Version = version
		config.Consumer.Return.Errors = true
	}

	// init consumer
	ctx := context.Background()
	client, err := sarama.NewConsumerGroup(kci.property.BrokerList, kci.property.GroupID, config)
	if err != nil {
		kci.error(err.Error())
		panic(err)
	}

	kci.consumerGroup = client
	defer func() {
		client.Close()
		if !kci.property.DisableAutoRestart {
			kci.info("the kafka consumer restarting")
			New(kci.property).Handle(hs, he)
		}
	}()

	started()
	defer func() {
		if err := recover(); err != nil {
			kci.error("the kafka consumer stopping because:", err)
		}
		closed()
	}()

	consumer := Consumer{
		kafkaConsumer: kci,
		handleSuccess: hs,
		handleError:   he,
	}
	go func() {
		for {
			consumer.ready = make(chan bool, 0)
			err := client.Consume(ctx, kci.property.Topics, &consumer)
			if err != nil {
				panic(err)
			}
		}
	}()
	<-consumer.ready
}
