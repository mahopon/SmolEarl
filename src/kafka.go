package main

import (
	"log"

	"github.com/IBM/sarama"
)

var KafkaProducer sarama.SyncProducer

func InitKafka() error {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	brokerList := []string{AppConfig.KafkaBroker}
	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return err
	}

	KafkaProducer = producer
	log.Println("Kafka producer initialized successfully")
	return nil
}

func CloseKafka() {
	if KafkaProducer != nil {
		KafkaProducer.Close()
	}
}

func PublishClickEvent(topic string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}

	_, _, err := KafkaProducer.SendMessage(msg)
	return err
}
