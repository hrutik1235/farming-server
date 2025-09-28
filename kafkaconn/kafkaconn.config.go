package kafkaconn

import (
	"context"
	"fmt"
	"net"

	"github.com/segmentio/kafka-go"
)

type KafkaConfig struct {
	brokers []string
	writer  *kafka.Writer
	readers []*kafka.Reader
}

func NewKafka(brokers []string) *KafkaConfig {
	return &KafkaConfig{
		brokers: brokers,
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:  brokers,
			Balancer: &kafka.LeastBytes{},
		}),
		readers: []*kafka.Reader{},
	}
}

func (k *KafkaConfig) CreateTopic(topic string) error {
	conn, err := kafka.Dial("tcp", k.brokers[0])
	if err != nil {
		return fmt.Errorf("failed to connect to broker: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("failed to get controller: %w", err)
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, fmt.Sprint(controller.Port)))
	if err != nil {
		return fmt.Errorf("failed to connect to controller: %w", err)
	}
	defer controllerConn.Close()

	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	return controllerConn.CreateTopics(topicConfig)
}

func (k *KafkaConfig) WriteMessage(ctx context.Context, topic string, message kafka.Message) error {
	k.writer.Topic = topic
	return k.writer.WriteMessages(ctx, message)
}

func (k *KafkaConfig) WriteMessages(ctx context.Context, topic string, messages ...kafka.Message) error {
	k.writer.Topic = topic
	return k.writer.WriteMessages(ctx, messages...)
}

func (k *KafkaConfig) NewReader(topic, groupID string) *kafka.Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  k.brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	k.readers = append(k.readers, reader)
	return reader
}

func (k *KafkaConfig) Consume(ctx context.Context, reader *kafka.Reader, handler func(kafka.Message)) error {
	for {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			return err
		}
		handler(m)
	}
}

func (k *KafkaConfig) Close() error {
	if err := k.writer.Close(); err != nil {
		return err
	}
	for _, r := range k.readers {
		if err := r.Close(); err != nil {
			return err
		}
	}
	return nil
}
