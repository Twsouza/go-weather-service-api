package services

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTemperatureByLocation(t *testing.T) {
	// Mock HTTP server
	mockResponse := `{"current": {"temp_c": 25.0}}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	service := NewWeatherAPIService("api-key", server.URL)

	temp, err := service.GetTemperatureByLocation("London")
	require.NoError(t, err)
	assert.Equal(t, 25.0, temp)
}

func TestCalculateTemperature(t *testing.T) {
	service := &WeatherAPIService{}

	fahrenheit, kelvin := service.CalculateTemperature(25.0)
	assert.Equal(t, 77.0, fahrenheit)
	assert.Equal(t, 298.15, kelvin)
}
