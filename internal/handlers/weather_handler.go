package handlers

import (
	"encoding/json"
	"net/http"

	"go-weather-service-api/internal/dto"
	"go-weather-service-api/internal/services"
)

const (
	ErrZipCodeNotFound = "can not find zipcode"
)

// WeatherHandler handles weather-related HTTP requests.
type WeatherHandler struct {
	ZipCodeService services.ZipCodeService
	WeatherService services.WeatherService
}

// ServeHTTP processes the HTTP request and returns the temperature data.
func (h *WeatherHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	if !h.ZipCodeService.IsValidCEP(cep) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	location, err := h.ZipCodeService.GetLocationByZipCode(cep)
	if err != nil {
		if err.Error() == ErrZipCodeNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	tempC, err := h.WeatherService.GetTemperatureByLocation(location)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tempF, tempK := h.WeatherService.CalculateTemperature(tempC)
	temp := dto.Temperature{
		Celsius:    tempC,
		Fahrenheit: tempF,
		Kelvin:     tempK,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(temp)
}
