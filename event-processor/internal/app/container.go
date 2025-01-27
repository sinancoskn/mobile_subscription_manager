package app

import (
	"event-processor/internal/config"
	"event-processor/internal/models"
	"event-processor/internal/services"

	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	container.Provide(config.LoadConfig)
	container.Provide(models.NewDB)
	container.Provide(models.NewWebhookRepository)
	container.Provide(services.NewWebhookService)
	container.Provide(models.NewManagerActionRepository)
	container.Provide(models.NewBatchRepository)
	container.Provide(services.NewWorkerManagerService)
	container.Provide(models.NewSubscriptionRepository)
	container.Provide(models.NewWorkerRepository)
	container.Provide(services.NewWorkerService)
	container.Provide(services.NewStoreApiService)

	return container
}
