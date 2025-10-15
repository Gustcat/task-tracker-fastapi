package kafka_common

import (
	"context"
	"github.com/segmentio/kafka-go"
	"time"
)

type Producer struct {
	writer *kafka.Writer
}

type ProducerConfig struct {
	BrokerAddrs []string
}

func NewProducer(cfg ProducerConfig) *Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.BrokerAddrs...),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		MaxAttempts:  5,
	}

	return &Producer{writer: writer}
}

func (p *Producer) Send(ctx context.Context, key, value []byte, topic string, time time.Time) error {
	msg := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time,
		Topic: topic,
	}

	return p.writer.WriteMessages(ctx, msg)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
