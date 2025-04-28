package kafka

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

var KafkaWriter *kafka.Writer

func InitKafka() {
	KafkaWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{os.Getenv("KAFKA_BROKER")}, // contoh: "localhost:9092"
		Topic:    "add-to-chart",
		Balancer: &kafka.LeastBytes{},
	})
}

func ProduceMessage(message []byte) error {
	err := KafkaWriter.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(time.Now().String()), // opsional, bisa pakai id transaksi
			Value: message,
		},
	)
	if err != nil {
		log.Println("Gagal kirim ke Kafka:", err)
		return err
	}
	return nil
}