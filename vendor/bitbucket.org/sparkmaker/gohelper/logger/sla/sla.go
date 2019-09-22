package sla

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/sparkmaker/gohelper/logger/stdout"

	"bitbucket.org/sparkmaker/gohelper/kafka/producer"
)

const (
	NA             = "n/a"
	XCorrelationID = "x-correlation-id"
	XActivity      = "x-activity"
	XAppID         = "x-app-id"
)

var (
	SLAKafkaTopic      = os.Getenv("SLA_KAFKA_TOPIC")
	SLAKafkaMaxRetry   = os.Getenv("SLA_KAFKA_MAX_RETRY")
	SLAKafkaKey        = os.Getenv("SLA_KAFKA_KEY")
	SLAKafkaBrokerList = strings.Split(os.Getenv("SLA_KAFKA_BROKER_LIST"), ",")
)

var (
	SLAServiceName  = os.Getenv("SLA_SERVICE_NAME")
	SLAServiceGroup = os.Getenv("SLA_SERVICE_GROUP")
	SLAVersion      = os.Getenv("SLA_VERSION")
	SLAInstanceID   = os.Getenv("SLA_INSTANCE_ID")
	SLASource       = os.Getenv("SLA_SOURCE")
)

type responseWriterRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *responseWriterRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

type message struct {
	InstanceID     string `json:"instanceId"`
	Source         string `json:"source"`
	HTTPPath       string `json:"httpPath"`
	HTTPResponse   int    `json:"httpResponse"`
	ServiceName    string `json:"serviceName"`
	ServiceGroup   string `json:"serviceGroup"`
	HTTPMethod     string `json:"httpMethod"`
	Version        string `json:"version"`
	CreatedTime    int64  `json:"createdTime"`
	ResponseTime   int64  `json:"responseTime"`
	XCorrelationID string `json:"correlationId"`
	XAppID         string `json:"appId"`
	XActivity      string `json:"activity"`
	SLA            int    `json:"sla"`
}

func (msg *message) fillIn() {
	if len(strings.TrimSpace(msg.InstanceID)) == 0 {
		msg.InstanceID = NA
	}
	if len(strings.TrimSpace(msg.Source)) == 0 {
		msg.Source = NA
	}
	if len(strings.TrimSpace(msg.HTTPPath)) == 0 {
		msg.HTTPPath = NA
	}
	if len(strings.TrimSpace(msg.ServiceName)) == 0 {
		msg.ServiceName = NA
	}
	if len(strings.TrimSpace(msg.ServiceGroup)) == 0 {
		msg.ServiceGroup = NA
	}
	if len(strings.TrimSpace(msg.Version)) == 0 {
		msg.Version = NA
	}
	if len(strings.TrimSpace(msg.HTTPMethod)) == 0 {
		msg.HTTPMethod = NA
	}
	if len(strings.TrimSpace(msg.XCorrelationID)) == 0 {
		msg.XCorrelationID = NA
	}
	if len(strings.TrimSpace(msg.XAppID)) == 0 {
		msg.XAppID = NA
	}
	if len(strings.TrimSpace(msg.XActivity)) == 0 {
		msg.XActivity = NA
	}
}

func (msg *message) calResponseTime() {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	msg.ResponseTime = now - msg.CreatedTime
}

func (msg *message) finished() {
	msg.fillIn()
	msg.calResponseTime()
}

type ServiceLevelAgreement struct {
	message                *message
	responseWriterRecorder *responseWriterRecorder
}

func pushMessageToKafka(msg *message) {
	defer func() { recover() }()
	slaKafkaMaxRetry, err := strconv.Atoi(SLAKafkaMaxRetry)
	if err != nil {
		slaKafkaMaxRetry = 0
	}
	producer := producer.KafkaProducer{
		Topic:      SLAKafkaTopic,
		MaxRetry:   slaKafkaMaxRetry,
		BrokerList: SLAKafkaBrokerList,
	}
	kResult := producer.Send(SLAKafkaKey, msg)
	bytes, _ := json.Marshal(msg)
	if kResult.Error != nil {
		stdout.Error("Error message:", kResult.Error.Error(), "Event message:", string(bytes))
		return
	}
	stdout.Info("Send a message to Kafka successfully:", string(bytes))
}

func (s *ServiceLevelAgreement) ResponseWriter() http.ResponseWriter {
	return s.responseWriterRecorder
}

func (s *ServiceLevelAgreement) Finished() {
	go func() {
		s.message.finished()
		s.message.HTTPResponse = s.responseWriterRecorder.status
		pushMessageToKafka(s.message)
	}()
}

func New(w http.ResponseWriter, r *http.Request) *ServiceLevelAgreement {
	header := map[string]string{}
	for k, v := range r.Header {
		header[strings.ToLower(strings.TrimSpace(k))] = strings.Join(v, ",")
	}
	// Initialize the status to 200 in case WriteHeader is not called
	return &ServiceLevelAgreement{&message{
		ServiceName:    SLAServiceName,
		ServiceGroup:   SLAServiceGroup,
		Version:        SLAVersion,
		InstanceID:     SLAInstanceID,
		Source:         SLASource,
		HTTPPath:       r.RequestURI,
		HTTPMethod:     r.Method,
		CreatedTime:    time.Now().UnixNano() / int64(time.Millisecond),
		XCorrelationID: header[XCorrelationID],
		XAppID:         header[XAppID],
		XActivity:      header[XActivity],
	}, &responseWriterRecorder{w, 200}}
}
