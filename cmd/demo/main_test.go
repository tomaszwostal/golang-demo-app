package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockDB is a mock type for gorm.DB
type MockDB struct {
	mock.Mock
}

func (m *MockDB) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

// Repository with MockDB
type TestRepository struct {
	DB *MockDB
}

func (r *TestRepository) CreatePlant(context *fiber.Ctx) error {
	plant := Plant{}

	if err := context.BodyParser(&plant); err != nil {
		return context.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "request failed",
		})
	}

	result := r.DB.Create(&plant)
	if result.Error != nil {
		return context.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "could not create plant",
		})
	}

	return context.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "plant created",
	})
}

func TestCreatePlant(t *testing.T) {
	mockDB := new(MockDB)
	repo := &TestRepository{DB: mockDB}

	app := fiber.New()
	app.Post("/api/create_plant", repo.CreatePlant)

	plant := Plant{Name: "Rose", Species: "Rosa", Plan: "Water daily"}
	body, _ := json.Marshal(plant)

	// Mock the Create method to return a gorm.DB with no error
	mockDB.On("Create", mock.AnythingOfType("*main.Plant")).Return(&gorm.DB{Error: nil})

	req := httptest.NewRequest(http.MethodPost, "/api/create_plant", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err, "App test should not return an error")
	assert.Equal(t, fiber.StatusOK, resp.StatusCode, "Status code should be 200 OK")

	// Check the response body
	var responseBody map[string]string
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err, "Should be able to decode response body")
	assert.Equal(t, "plant created", responseBody["message"], "Response message should be 'plant created'")

	mockDB.AssertExpectations(t)
}

func TestCreatePlantError(t *testing.T) {
	mockDB := new(MockDB)
	repo := &TestRepository{DB: mockDB}

	app := fiber.New()
	app.Post("/api/create_plant", repo.CreatePlant)

	plant := Plant{Name: "Rose", Species: "Rosa", Plan: "Water daily"}
	body, _ := json.Marshal(plant)

	// Mock the Create method to return a gorm.DB with an error
	mockDB.On("Create", mock.AnythingOfType("*main.Plant")).Return(&gorm.DB{Error: gorm.ErrInvalidData})

	req := httptest.NewRequest(http.MethodPost, "/api/create_plant", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err, "App test should not return an error")
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode, "Status code should be 400 Bad Request")

	// Check the response body
	var responseBody map[string]string
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err, "Should be able to decode response body")
	assert.Equal(t, "could not create plant", responseBody["message"], "Response message should be 'could not create plant'")

	mockDB.AssertExpectations(t)
}
