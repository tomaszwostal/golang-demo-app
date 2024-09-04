package storage // Ensure this matches your main file's package name

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigString(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		expected string
	}{
		{
			name: "Full config",
			config: &Config{
				Host:     "localhost",
				Port:     "5432",
				User:     "user",
				Password: "password",
				DBName:   "testdb",
				SSLMode:  "disable",
			},
			expected: "host=localhost port=5432 user=user dbname=testdb password=password sslmode=disable",
		},
		{
			name: "Partial config",
			config: &Config{
				Host:     "localhost",
				Port:     "5432",
				User:     "user",
				Password: "",
				DBName:   "testdb",
				SSLMode:  "",
			},
			expected: "host=localhost port=5432 user=user dbname=testdb password= sslmode=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewConnection(t *testing.T) {
	// This is a placeholder test. It doesn't actually test the connection,
	// but ensures that the function exists and can be called.
	config := &Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "user",
		Password: "password",
		DBName:   "testdb",
		SSLMode:  "disable",
	}

	// We're not actually opening a connection here
	_, err := NewConnection(config)

	// We're just checking that the function runs without panicking
	// In a real scenario, this would likely return an error because we're not actually connecting to a database
	assert.Error(t, err, "Expected an error when trying to connect to a non-existent database")
}
