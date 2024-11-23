package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

// WeatherService defines the interface for fetching temperature by location.
type WeatherService interface {
	GetTemperatureByLocation(location string) (float64, error)
	CalculateTemperature(tempC float64) (fahrenheit float64, kelvin float64)
}

// WeatherAPIService implements WeatherService using the WeatherAPI.
type WeatherAPIService struct {
	ApiKey  string
	BaseURL string
}

func NewWeatherAPIService(apiKey string, baseURL string) *WeatherAPIService {
	return &WeatherAPIService{
		ApiKey:  apiKey,
		BaseURL: baseURL,
	}
}

// weatherAPIResponse represents the response from the WeatherAPI.
type weatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

// GetTemperatureByLocation fetches the current temperature for a given location.
func (w *WeatherAPIService) GetTemperatureByLocation(location string) (float64, error) {
	apiURL, err := url.Parse(w.BaseURL)
	if err != nil {
		return 0, err
	}
	apiURL.Path = "/v1/current.json"

	query := apiURL.Query()
	query.Set("key", w.ApiKey)
	query.Set("q", location)
	query.Set("aqi", "no")
	apiURL.RawQuery = query.Encode()

	logrus.WithFields(logrus.Fields{
		"apiURL": apiURL.String(),
	}).Info("Fetching temperature by location", location)

	resp, err := http.Get(apiURL.String())
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("failed to get weather data")
	}

	var data weatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	logrus.WithFields(logrus.Fields{
		"tempC": data.Current.TempC,
	}).Info("Temperature found for location", location)

	return data.Current.TempC, nil
}

// CalculateTemperature converts the temperature from Celsius to Fahrenheit and Kelvin.
func (w *WeatherAPIService) CalculateTemperature(tempC float64) (float64, float64) {
	return tempC*1.8 + 32, tempC + 273.15
}
