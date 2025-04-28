package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"golang_api/model" // Update dengan path yang sesuai untuk model kamu

	"github.com/segmentio/kafka-go"
)

// StartConsumer akan memulai Kafka consumer dan menangani pesan yang diterima
func StartConsumer() {
	// Inisialisasi Kafka Consumer
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("KAFKA_BROKER")}, // Kafka Broker
		Topic:   "add-to-chart",                      // Nama Topic yang akan di-consume
		GroupID: "chart-consumer-group",               // Group ID untuk consumer
		MinBytes: 10e3,  // 10KB
		MaxBytes: 10e6,  // 10MB
	})

	defer consumer.Close()

	log.Println("Kafka Consumer started. Waiting for messages...")

	// Loop untuk membaca pesan dari Kafka
	for {
		// Membaca pesan dari Kafka
		message, err := consumer.FetchMessage(context.Background())
		if err != nil {
			log.Println("Error membaca pesan dari Kafka:", err)
			continue
		}

		// Menampilkan informasi tentang pesan yang diterima
		log.Printf("Pesan diterima: %s\n", string(message.Value))
		

		// Decode pesan JSON ke dalam struct ChartItem
		var chartRequest model.ChartItem
		err = json.Unmarshal(message.Value, &chartRequest)
		if err != nil {
			log.Println("Error mendecode pesan JSON:", err)
			continue
		}

		// Menambahkan timestamp untuk created_at dan updated_at
		createdAt := time.Now().Unix()
		updatedAt := createdAt

		// Simpan data ke database
		_, err = model.CreateChart(chartRequest.UserID, chartRequest.ProductID, chartRequest.Quantity, int(createdAt), int(updatedAt))
		if err != nil {
			log.Println("Error menyimpan ke database:", err)
		} else {
			log.Printf("Data berhasil disimpan ke database: %+v\n", chartRequest)
		}

		log.Printf("Menerima data untuk user_id: %d, product_id: %d, quantity: %d", chartRequest.UserID, chartRequest.ProductID, chartRequest.Quantity)


		// Kalau sukses baru commit
		if err := consumer.CommitMessages(context.Background(), message); err != nil {
			log.Println("Error commit message:", err)
		} else {
			log.Printf("Data berhasil disimpan: %+v", chartRequest)
		}
	}
}
