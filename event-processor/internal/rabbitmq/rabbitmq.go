package rabbitmq

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func ConnectRabbitMQ(rabbitMQURL string) (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	return conn, ch
}

func ConsumeMessages(ch *amqp.Channel, queueName string, handler func(ctx context.Context, data map[string]interface{}) error) {
	msgs, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var msg map[string]interface{}

			err := json.Unmarshal(d.Body, &msg)
			if err != nil {
				log.Printf("Error decoding message: %v", err)
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			err = handler(ctx, msg)
			if err != nil {
				log.Printf("Error handling message: %v", err)
			}
		}
	}()

	log.Println("Waiting for messages...")
	<-forever
}
