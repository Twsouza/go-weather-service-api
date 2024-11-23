package main

import (
	"log"
	"net/http"
	"os"

	"go-weather-service-api/internal/handlers"
	"go-weather-service-api/internal/services"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("Could not load .env file")
	}

	if os.Getenv("WEATHER_API_KEY") == "" {
		log.Fatal("WEATHER_API_KEY must be set")
	}

	zipCodeService := services.NewViaCEPService(
		"https://viacep.com.br",
	)
	weatherService := services.NewWeatherAPIService(
		os.Getenv("WEATHER_API_KEY"),
		"http://api.weatherapi.com",
	)

	weatherHandler := &handlers.WeatherHandler{
		ZipCodeService: zipCodeService,
		WeatherService: weatherService,
	}

	http.Handle("/weather", weatherHandler)
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
