package kafka

import (
	"context"
	"log"
	"strings"

	"github.com/IBM/sarama"
)

// StartEventsConsumer starts a simple Kafka consumer that listens to the given topic
func StartEventsConsumer(ctx context.Context, brokersCSV, topic string, OnEvent func(p int32, pc sarama.PartitionConsumer)) error {
	if brokersCSV == "" {
		log.Println("Kafka brokers are empty; skipping events consumer")
		return nil
	}

	brokers := strings.Split(brokersCSV, ",")

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return err
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Printf("Error closing Kafka consumer: %v", err)
		}
	}()

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		return err
	}

	log.Printf("Payment Kafka consumer starting; topic=%s partitions=%v", topic, partitions)

	done := make(chan struct{})
	for _, partition := range partitions {
		pc, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Printf("Error starting partition consumer p=%d: %v", partition, err)
			continue
		}

		go OnEvent(partition, pc)
	}

	<-ctx.Done()
	close(done)
	return nil
}
