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

type Key int

const (
	ID Key = iota
)

type Consumer struct {
	Ready chan struct{}
	wg    *sync.WaitGroup
	pool  *ants.Pool
	svc   *Service
}

func NewConsumer(wg *sync.WaitGroup, pool *ants.Pool, svc *Service) *Consumer {
	return &Consumer{
		Ready: make(chan struct{}),
		wg:    wg,
		pool:  pool,
		svc:   svc,
	}
}

func (h *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(h.Ready)
	return nil
}

func (h *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		msg := message
		ctx := context.WithValue(context.Background(), ID, uuid.New().String())

		log.Printf("Receive message : value : %v , timestamp : %v , topic : %v", string(message.Value), message.Timestamp, message.Topic)

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

			log.Printf("Consumer success Mark message : value : %v , timestamp : %v , topic : %v", string(message.Value), message.Timestamp, message.Topic)

		})
	}

	return nil
}
