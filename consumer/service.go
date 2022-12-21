package consumer

import (
	"context"
	"encoding/json"
	"kafka-reconsign/internal/util"
	"kafka-reconsign/model"
	"log"

	"github.com/Shopify/sarama"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

// {"refId": "refId","state": "state","message": "message"}
func (s *Service) Process(ctx context.Context, kafkaHeader []*sarama.RecordHeader, rawMessage string) error {
	header := populateMessageKafkaHeader(ctx, kafkaHeader)

	var kafkaMsg model.KafkaModel
	err := json.Unmarshal([]byte(rawMessage), &kafkaMsg)
	if err != nil {
		return err
	}

	log.Println("Get Kafka Header: ", util.StructToString(header))
	log.Println("Get Kafka Message: ", util.StructToString(kafkaMsg))

	return nil

}

func populateMessageKafkaHeader(ctx context.Context, headers []*sarama.RecordHeader) model.Header {
	var h model.Header
	// for _, header := range headers {
	// 	if string(header.Key) != "" && strings.EqualFold(string(header.Key), modelx.HeaderXRefID) {
	// 		h.ReferenceId = string(header.Value)
	// 	}
	// }
	return h
}
