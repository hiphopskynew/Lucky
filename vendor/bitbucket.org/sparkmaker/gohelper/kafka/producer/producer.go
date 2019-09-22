package producer

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
)

type KafkaResult struct {
	Topic     string      `json:"topic"`
	Key       string      `json:"key"`
	Message   interface{} `json:"message"`
	Partition int32       `json:"partition,omitempty"`
	Offset    int64       `json:"offset,omitempty"`
	Error     error       `json:"error,omitempty"`
}

type KafkaProducer struct {
	BrokerList             []string
	Topic                  string
	MaxRetry               int
	maxSizeBytesPerMessage int64
}

var defaultKafkaMessageSize int64 = 950000

func (kafka *KafkaProducer) SetBrokerList(broker []string) {
	kafka.BrokerList = broker
}

func (kafka *KafkaProducer) SetTopic(topic string) {
	kafka.Topic = topic
}

func (kafka *KafkaProducer) SetMaxRetry(retry int) {
	kafka.MaxRetry = retry
}

func (kafka *KafkaProducer) SetMaxSizeMessage(size int64) error {
	if size <= 0 || size > defaultKafkaMessageSize {
		panic(errors.New("Size must be more than 0 and less than " + strconv.FormatInt(defaultKafkaMessageSize, 10)))
	}
	kafka.maxSizeBytesPerMessage = size
	return nil
}

func (kafka *KafkaProducer) sends(key string, message interface{}, wg *sync.WaitGroup) KafkaResult {
	defer wg.Done()
	return kafka.Send(key, message)
}

func (kafka *KafkaProducer) Send(key string, message interface{}) KafkaResult {
	result := KafkaResult{
		Topic:   kafka.Topic,
		Key:     key,
		Message: message,
		Error:   nil,
	}
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = kafka.MaxRetry
	config.Producer.Return.Successes = true

	var msgArrBytes []byte
	switch message.(type) {
	case string:
		msgArrBytes = []byte(message.(string))
	default:
		b, err := json.Marshal(message)
		if err != nil {
			result.Error = err
			return result
		}
		msgArrBytes = b
	}

	producer, err := sarama.NewSyncProducer(kafka.BrokerList, config)
	if err != nil {
		result.Error = err
		return result
	}

	defer func() {
		if err := producer.Close(); err != nil {
			panic(err)
		}
	}()

	var encoderKey sarama.Encoder
	if len(strings.TrimSpace(key)) != 0 {
		encoderKey = sarama.StringEncoder(key)
	}

	partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic: kafka.Topic,
		Key:   encoderKey,
		Value: sarama.StringEncoder(string(msgArrBytes)),
	})

	if err != nil {
		result.Error = err
		return result
	}

	result.Partition = partition
	result.Offset = offset
	result.Error = err

	return result
}

func (kafka *KafkaProducer) SendSegment(key string, message string) error {
	if kafka.maxSizeBytesPerMessage == 0 {
		kafka.maxSizeBytesPerMessage = defaultKafkaMessageSize // default maxSizeBytePerMessage
	}

	wg := &sync.WaitGroup{}
	msgs := splitMsg(message, int(kafka.maxSizeBytesPerMessage))
	msgsLen := int64(len(msgs))
	sid := uuid.New().String()
	for i, msg := range msgs {
		segment := segment{
			SID:      sid,
			SOrder:   int64(i + 1),
			SOverall: msgsLen,
			SData:    msg,
		}
		wg.Add(1)
		go kafka.sends(key, segment, wg)
	}
	wg.Wait()
	return nil
}

func splitMsg(msg string, length int) []string {
	start := 0
	end := length
	msgByte := []byte(msg)
	grouped := []string{}
	for end < len(msgByte) {
		grouped = append(grouped, string(msgByte[start:end]))
		start = end
		end = end + length
	}
	grouped = append(grouped, string(msgByte[start:len(msgByte)]))
	return grouped
}

func (kafka *KafkaProducer) produceSegmentMessage(wg *sync.WaitGroup, producer sarama.SyncProducer, segment segment, key string) {
	defer wg.Done()
	b, _ := json.Marshal(&segment)
	_, _, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic: kafka.Topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(string(b)),
	})
	if err != nil {
		panic(err.Error())
	}
}

type segment struct {
	SID      string `json:"s_id"`
	SOrder   int64  `json:"s_order"`
	SOverall int64  `json:"s_overall"`
	SData    string `json:"s_data"`
}
