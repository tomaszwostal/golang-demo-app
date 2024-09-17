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

type Plant struct {
    Name    string `json:"name"`
    Species string `json:"species"`
    Plan    string `json:"plan"`
}

func (r *Repository) CreatePlant(context *fiber.Ctx) error {
    plant := Plant{}

    err := context.BodyParser(&plant)
    if err != nil {
        context.Status(http.StatusUnprocessableEntity).JSON(
            &fiber.Map{"message": "request failed"})
        return err
    }

    err = r.DB.Create(&plant).Error
    if err != nil {
        context.Status(http.StatusBadRequest).JSON(
            &fiber.Map{"message": "could not create plant"})
        return err
    }

    context.Status(http.StatusOK).JSON(&fiber.Map{
        "message": "plant created"})
    return nil
}

func (r *Repository) GetPlants(context *fiber.Ctx) error {
    plantModels := &[]models.Plants{}

    err := r.DB.Find(plantModels).Error
    if err != nil {
        context.Status(http.StatusBadRequest).JSON(
            &fiber.Map{"message": "could not get plants"})
        return err
    }

    context.Status(http.StatusOK).JSON(&fiber.Map{
        "message": "plants fetched successfully",
        "data":    plantModels,
    })
    return nil
}

func (r *Repository) DeletePlant(context *fiber.Ctx) error {
    plantModel := models.Plants{}
    id := context.Params("id")

    if id == "" {
        context.Status(http.StatusBadRequest).JSON(
            &fiber.Map{"message": "id is required"})
        return nil
    }
    err := r.DB.Delete(&plantModel, id).Error
    if err != nil {
        context.Status(http.StatusBadRequest).JSON(
            &fiber.Map{"message": "could not delete plant"})
        return err
    }

    context.Status(http.StatusOK).JSON(&fiber.Map{
        "message": "plant deleted successfully"})
    return nil
}

func (r *Repository) GetPlantByID(context *fiber.Ctx) error {
    plantModel := &models.Plants{}
    id := context.Params("id")

    if id == "" {
        context.Status(http.StatusBadRequest).JSON(
            &fiber.Map{"message": "id is required"})
        return nil
    }
    fmt.Println("The id is: ", id)

    err := r.DB.Where("id = ?", id).First(plantModel).Error
    if err != nil {
        context.Status(http.StatusBadRequest).JSON(
            &fiber.Map{"message": "could not get plant"})
        return err
    }
    context.Status(http.StatusOK).JSON(&fiber.Map{
        "message": "plant fetched successfully",
        "data":    plantModel,
    })

    return nil
}

func (r *Repository) UpdatePlant(context *fiber.Ctx) error {
    plantModel := models.Plants{}
    id := context.Params("id")

    if id == "" {
        context.Status(http.StatusBadRequest).JSON(
            &fiber.Map{"message": "id is required"})
        return nil
    }

    err := r.DB.Where("id = ?", id).First(&plantModel).Error
    if err != nil {
        context.Status(http.StatusBadRequest).JSON(
            &fiber.Map{"message": "could not get plant"})
        return err
    }
    err = context.BodyParser(&plantModel)
    if err != nil {
        context.Status(http.StatusUnprocessableEntity).JSON(
            &fiber.Map{"message": "request failed"})
        return err
    }
    err = r.DB.Save(&plantModel).Error
    if err != nil {
        context.Status(http.StatusBadRequest).JSON(
            &fiber.Map{"message": "could not update plant"})
        return err
    }
    context.Status(http.StatusOK).JSON(&fiber.Map{
        "message": "plant updated successfully",
        "data":    plantModel,
    })
    return nil
}

// Dodajemy brakującą metodę SetupRoutes
func (r *Repository) SetupRoutes(app *fiber.App) {
    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })
    api := app.Group("/api")
    api.Post("/create_plant", r.CreatePlant)
    api.Delete("/delete_plant/:id", r.DeletePlant)
    api.Get("/get_plant/:id", r.GetPlantByID)
    api.Get("/get_plants", r.GetPlants)
    api.Put("/update_plant/:id", r.UpdatePlant)
}

func InitOpenTelemetry(ctx context.Context) (*trace.TracerProvider, error) {
    collectorEndpoint := os.Getenv("OTEL_COLLECTOR_ENDPOINT")
    if collectorEndpoint == "" {
        collectorEndpoint = "localhost:4317" // Domyślny adres collectora
    }

    exporter, err := otlptracegrpc.New(ctx,
        otlptracegrpc.WithEndpoint(collectorEndpoint),
        otlptracegrpc.WithInsecure(), // Użyj WithInsecure(), jeśli collector nie używa TLS
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

    // err = godotenv.Load()
    // if err != nil {
    //     log.Fatal("Error loading environment variables")
    // }

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

