package main

import (
    "context"
    "fmt"
    "github.com/gofiber/contrib/otelfiber"
    "github.com/gofiber/fiber/v2"
    "github.com/joho/godotenv"
    "github.com/tomaszwostal/golang-demo-app/models"
    "github.com/tomaszwostal/golang-demo-app/storage"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/sdk/resource"
    "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
    "gorm.io/gorm"
    "log"
    "net/http"
    "os"
)

type Repository struct {
    DB *gorm.DB
}

// ... (pozostała część kodu)

func InitOpenTelemetry(ctx context.Context) (*trace.TracerProvider, error) {
    collectorEndpoint := os.Getenv("OTEL_COLLECTOR_ENDPOINT")
    if collectorEndpoint == "" {
        collectorEndpoint = "localhost:4317"
    }

    exporter, err := otlptracegrpc.New(ctx,
        otlptracegrpc.WithEndpoint(collectorEndpoint),
        otlptracegrpc.WithInsecure(),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create the collector exporter: %w", err)
    }

    resource, err := resource.New(ctx,
        resource.WithAttributes(
            semconv.ServiceNameKey.String("golang-demo-app"),
        ),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create resource: %w", err)
    }

    tracerProvider := trace.NewTracerProvider(
        trace.WithBatcher(exporter),
        trace.WithResource(resource),
    )

    otel.SetTracerProvider(tracerProvider)

    return tracerProvider, nil
}

func main() {
    ctx := context.Background()
    tp, err := InitOpenTelemetry(ctx)
    if err != nil {
        log.Fatalf("failed to initialize OpenTelemetry: %v", err)
    }
    defer func() {
        if err := tp.Shutdown(ctx); err != nil {
            log.Fatalf("Error shutting down tracer provider: %v", err)
        }
    }()

    err = godotenv.Load()
    if err != nil {
        log.Fatal("Error loading environment variables")
    }

    config := &storage.Config{
        Host:     os.Getenv("DB_HOST"),
        Port:     os.Getenv("DB_PORT"),
        User:     os.Getenv("DB_USER"),
        Password: os.Getenv("DB_PASSWORD"),
        DBName:   os.Getenv("DB_NAME"),
        SSLMode:  os.Getenv("DB_SSLMODE"),
    }

    db, err := storage.NewConnection(config)
    if err != nil {
        log.Fatal("Error connecting to database")
    }

    err = models.MigratePlants(db)
    if err != nil {
        log.Fatal("Error migrating plants table")
    }

    r := Repository{
        DB: db,
    }

    app := fiber.New()

    // Dodajemy middleware OpenTelemetry do Fiber
    app.Use(otelfiber.Middleware())

    r.SetupRoutes(app)
    app.Listen(":3000")
}

