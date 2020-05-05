package main

import (
	"log"
	"os"
	"time"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func main() {
	broker := os.Getenv("BROKER")
	topic := os.Getenv("TOPIC")

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})

	if err != nil {
		log.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	log.Printf("Created Producer %v\n", p)

	for {
		deliveryChan := make(chan kafka.Event)

		value := "Hello from k8s"
		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(value),
			Headers:        []kafka.Header{{Key: "TestHeader", Value: []byte("header values are binary")}},
		}, deliveryChan)
		if err != nil {
			log.Println(err.Error())
		}

		e := <-deliveryChan
		m := e.(*kafka.Message)

		if m.TopicPartition.Error != nil {
			log.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
		} else {
			log.Printf("Delivered message to topic %s [%d] at offset %v\n",
				*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
		}
		close(deliveryChan)
		time.Sleep(time.Second * 3)
	}
}
