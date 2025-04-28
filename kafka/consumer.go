package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"golang_api/model"

	"github.com/segmentio/kafka-go"
)

// StartConsumer akan memulai Kafka consumer dan menangani pesan yang diterima
func StartConsumer() {
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("KAFKA_BROKER")},
		Topic:   "add-to-chart",
		GroupID: "chart-consumer-group",
		MinBytes: 10e3,  // 10KB
		MaxBytes: 10e6,  // 10MB
	})

	defer consumer.Close()

	log.Println("Kafka Consumer started. Waiting for messages...")

	for {
		message, err := consumer.FetchMessage(context.Background())
		if err != nil {
			log.Println("Error membaca pesan dari Kafka:", err)
			continue
		}

		log.Printf("Pesan diterima: %s\n", string(message.Value))

		var chartRequest model.ChartItem
		err = json.Unmarshal(message.Value, &chartRequest)
		if err != nil {
			log.Println("Error mendecode pesan JSON:", err)
			continue
		}

		// Simpan data chart ke database
		chart, err := model.CreateChart(chartRequest.UserID, chartRequest.ProductID, chartRequest.Quantity)
		if err != nil {
			log.Println("Error menyimpan ke database:", err)
		} else {
			log.Printf("Data berhasil disimpan ke database: %+v\n", chart)
			
			// Catat ke tabel chart_event setelah berhasil simpan ke chart
			action := "insert"
			if chartRequest.Quantity > 1 {
				action = "update"
			}

			err = model.CreateChartEvent(chartRequest.UserID, chartRequest.ProductID, chartRequest.Quantity, action)
			if err != nil {
				log.Println("Error menyimpan event ke chart_event:", err)
			} else {
				log.Printf("Event berhasil dicatat: user_id=%d, product_id=%d, action=%s", chartRequest.UserID, chartRequest.ProductID, action)
			}
		}

		if err := consumer.CommitMessages(context.Background(), message); err != nil {
			log.Println("Error commit message:", err)
		} else {
			log.Printf("Pesan berhasil di-commit")
		}
	}
}

