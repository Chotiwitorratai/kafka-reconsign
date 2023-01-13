package consumer

import (
	"context"
	"log"
	"runtime/debug"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"github.com/panjf2000/ants/v2"
)

type Key2 int

const (
	ID2 Key = iota
)

type Consumer2 struct {
	Ready chan struct{}
	wg    *sync.WaitGroup
	pool  *ants.Pool
	svc   ServiceInsurance
}

func NewConsumer2(wg *sync.WaitGroup, pool *ants.Pool, svc ServiceInsurance) *Consumer2 {
	return &Consumer2{
		Ready: make(chan struct{}),
		wg:    wg,
		pool:  pool,
		svc:   svc,
	}
}

func (h *Consumer2) Setup(sarama.ConsumerGroupSession) error {
	close(h.Ready)
	return nil
}

func (h *Consumer2) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Consumer2) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		msg := message
		ctx := context.WithValue(context.Background(), ID, uuid.New().String())

		log.Printf("Receive2 message : value : %v , timestamp : %v , topic : %v", string(message.Value), message.Timestamp, message.Topic)

		session.MarkMessage(msg, "")

		h.wg.Add(1)
		_ = h.pool.Submit(func() {
			defer func() {
				if r := recover(); r != nil {
					log.Println(string(debug.Stack()))
				}
			}()

			defer h.wg.Done()

			if err := h.svc.Process(msg.Topic, ctx, msg.Headers, string(msg.Value)); err != nil {
				log.Println("Process Err: ", err)
				return
			}

			log.Printf("Consumer2 success Mark message : value : %v , timestamp : %v , topic : %v", string(message.Value), message.Timestamp, message.Topic)

		})
	}

	return nil
}
