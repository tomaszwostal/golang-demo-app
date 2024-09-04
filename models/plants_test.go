package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func TestMigratePlants(t *testing.T) {
	// Create an in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Run the migration
	err = MigratePlants(db)
	if err != nil {
		t.Fatalf("Failed to migrate plants: %v", err)
	}

	// Check if the table was created
	if !db.Migrator().HasTable(&Plants{}) {
		t.Fatalf("Expected table Plants to be created")
	}
}

func TestCreatePlant(t *testing.T) {
	// Create an in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Run the migration
	err = MigratePlants(db)
	if err != nil {
		t.Fatalf("Failed to migrate plants: %v", err)
	}

	// Create a new plant
	plant := Plants{
		Name:    stringPtr("Rose"),
		Species: stringPtr("Rosa"),
		Plan:    stringPtr("First"),
	}
	result := db.Create(&plant)
	if result.Error != nil {
		t.Fatalf("Failed to create plant: %v", result.Error)
	}

	// Check if the plant was created
	var count int64
	db.Model(&Plants{}).Count(&count)
	if count != 1 {
		t.Fatalf("Expected 1 plant, got %d", count)
	}
}

func stringPtr(s string) *string {
	return &s
}
