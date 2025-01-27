package app

import (
	"log"

	"event-processor/internal/config"
	"event-processor/internal/rabbitmq"
	"event-processor/internal/services"
)

func StartCallbackApp() error {
	container := BuildContainer()

	err := container.Invoke(func(config *config.Config, service *services.WebhookService) {
		conn, ch := rabbitmq.ConnectRabbitMQ(config.RabbitMQURL)
		defer conn.Close()
		defer ch.Close()

		// Queue to consume messages from
		queue := "subscription_events"
		log.Printf("Listening to queue: %s", queue)

		// Start consuming messages
		rabbitmq.ConsumeMessages(ch, queue, service.HandleEvent)
	})

	return err
}
