package workermanager

import (
	"context"
	"event-processor/internal/config"
	"event-processor/internal/models"
	"event-processor/internal/services"
	"fmt"
	"log"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/dig"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	Container   *dig.Container
	pgContainer testcontainers.Container
)

func SetupTest() {
	Container = dig.New()

	// Start a PostgreSQL test container
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15", // Use the desired PostgreSQL version
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	var err error
	pgContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Failed to start PostgreSQL container: %v", err)
	}

	// Get the container's host and port
	host, err := pgContainer.Host(ctx)
	if err != nil {
		log.Fatalf("Failed to get container host: %v", err)
	}
	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("Failed to get container port: %v", err)
	}

	// Connect to the PostgreSQL database
	dsn := fmt.Sprintf("host=%s port=%s user=test password=test dbname=testdb sslmode=disable", host, port.Port())
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL test database: %v", err)
	}

	// Auto-migrate the database schema
	err = db.AutoMigrate(&models.Subscription{}, &models.ManagerAction{}, &models.Batch{}, &models.Subscription{}, &models.Webhook{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}

	// Provide dependencies to the DI container
	Container.Provide(func() *gorm.DB {
		return db
	})

	Container.Provide(config.LoadConfig)
	Container.Provide(models.NewWebhookRepository)
	Container.Provide(services.NewWebhookService)
	Container.Provide(models.NewManagerActionRepository)
	Container.Provide(models.NewBatchRepository)
	Container.Provide(services.NewWorkerManagerService)
	Container.Provide(models.NewSubscriptionRepository)
	Container.Provide(models.NewWorkerRepository)
	Container.Provide(services.NewWorkerService)
	Container.Provide(services.NewStoreApiService)
}

func TearDownTest() {
	if pgContainer != nil {
		_ = pgContainer.Terminate(context.Background())
	}
}
