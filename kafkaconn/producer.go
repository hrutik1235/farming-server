package kafkaconn

import (
	"context"
	"fmt"

	"net"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	brokers []string
}

func (k *KafkaProducer) Produce(topic string, broker string) error {
	conn, err := kafka.Dial("tcp", broker)

	if err != nil {
		return err
	}

	defer conn.Close()

	controller, err := conn.Controller()

	if err != nil {
		return err
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, fmt.Sprint(controller.Port)))

	if err != nil {
		return err
	}

	defer controllerConn.Close()

	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	return controllerConn.CreateTopics(topicConfig)
}

func (k *KafkaProducer) Write(topic string, message kafka.Message) error {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  k.brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})

	ctx := context.Background()

	return writer.WriteMessages(ctx, message)
}
