package handlers

import (
	"encoding/json"
	"errors"
	"go-weather-service-api/internal/dto"
	"go-weather-service-api/mocks/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeatherHandlerServeHTTP(t *testing.T) {
	var celsius float64 = 25
	var fahrenheit float64 = 77
	var kelvin float64 = 298.15

	t.Run("Valid zipcode", func(t *testing.T) {
		zipMock := services.NewMockZipCodeService(t)
		zipMock.EXPECT().GetLocationByZipCode("12345678").Return("São Paulo", nil)
		zipMock.EXPECT().IsValidCEP("12345678").Return(true)
		weatherMock := services.NewMockWeatherService(t)
		weatherMock.EXPECT().GetTemperatureByLocation("São Paulo").Return(celsius, nil)
		weatherMock.EXPECT().CalculateTemperature(celsius).Return(fahrenheit, kelvin)
		weatherHandler := WeatherHandler{
			ZipCodeService: zipMock,
			WeatherService: weatherMock,
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/weather?cep=12345678", nil)
		weatherHandler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		result := dto.Temperature{}
		assert.NoError(t, json.NewDecoder(w.Body).Decode(&result))
		assert.Equal(t, celsius, result.Celsius)
		assert.Equal(t, kelvin, result.Kelvin)
	})

	t.Run("Invalid zipcode", func(t *testing.T) {
		zipMock := services.NewMockZipCodeService(t)
		zipMock.EXPECT().IsValidCEP("12345678").Return(false)
		weatherHandler := WeatherHandler{
			ZipCodeService: zipMock,
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/weather?cep=12345678", nil)
		weatherHandler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("Zipcode not found", func(t *testing.T) {
		zipMock := services.NewMockZipCodeService(t)
		zipMock.EXPECT().GetLocationByZipCode("12345678").Return("", errors.New(ErrZipCodeNotFound))
		zipMock.EXPECT().IsValidCEP("12345678").Return(true)
		weatherHandler := WeatherHandler{
			ZipCodeService: zipMock,
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/weather?cep=12345678", nil)
		weatherHandler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Zip code services returns an error", func(t *testing.T) {
		zipMock := services.NewMockZipCodeService(t)
		zipMock.EXPECT().GetLocationByZipCode("12345678").Return("", errors.New("error"))
		zipMock.EXPECT().IsValidCEP("12345678").Return(true)
		weatherHandler := WeatherHandler{
			ZipCodeService: zipMock,
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/weather?cep=12345678", nil)
		weatherHandler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Weather service returns an error", func(t *testing.T) {
		zipMock := services.NewMockZipCodeService(t)
		zipMock.EXPECT().GetLocationByZipCode("12345678").Return("São Paulo", nil)
		zipMock.EXPECT().IsValidCEP("12345678").Return(true)
		weatherMock := services.NewMockWeatherService(t)
		weatherMock.EXPECT().GetTemperatureByLocation("São Paulo").Return(0, errors.New("error"))
		weatherHandler := WeatherHandler{
			ZipCodeService: zipMock,
			WeatherService: weatherMock,
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/weather?cep=12345678", nil)
		weatherHandler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Invalid request", func(t *testing.T) {
		zipMock := services.NewMockZipCodeService(t)
		zipMock.EXPECT().IsValidCEP("").Return(false)
		weatherHandler := WeatherHandler{
			ZipCodeService: zipMock,
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/weather", nil)
		weatherHandler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
}
