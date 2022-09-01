package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	config "gitlab.ozon.dev/pircuser61/catalog_iface/config"
)

type Consumer struct{}

func (c *Consumer) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case <-session.Context().Done():
			fmt.Println("Error's session context done")
			return nil
		case msg, ok := <-claim.Messages():
			if !ok {
				fmt.Println("Error's channel closed")
				return nil
			}
			fmt.Printf("Произошла ошибка: %v data: %v\n", msg.Timestamp, string(msg.Value))
			session.MarkMessage(msg, "")
		}
	}
}

func runErrListiner(ctx context.Context) {
	/*
		читает вообще все ошибки, независимо относятся они к данному клиенту или нет
		и выводит на экран

	*/

	var topics = []string{config.Topic_error}

	kafkaCfg := sarama.NewConfig()
	kafkaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewConsumerGroup(config.KafkaBrokers, "catalog", kafkaCfg)
	if err != nil {
		fmt.Printf("Cant create consumer: %v\n", err)
		return
	}

	consumer := &Consumer{}
	for {
		err := client.Consume(ctx, topics, consumer)
		if err != nil {
			fmt.Printf("Err on consume: %v", err)
			time.Sleep(time.Second * 10)
		}
	}
}
