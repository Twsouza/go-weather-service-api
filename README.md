# Weather Application in Go

## Table of Contents

- [Description](#description)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Testing](#testing)
- [Dockerization](#dockerization)
- [Deployment on Google Cloud Run](#deployment-on-google-cloud-run)
- [Contributing](#contributing)
- [License](#license)


## Description

This Go-based web application receives a Brazilian postal code (CEP), identifies the corresponding city, and returns the current temperature in Celsius, Fahrenheit, and Kelvin. It utilizes the ViaCEP API to find the location and the WeatherAPI to get the current temperature.

The application is available through the Cloud Run in the following URL: https://weather-app-16583978889.us-central1.run.app/weather?cep=05436060


## Features

- **CEP Validation**: Accepts and validates an 8-digit Brazilian CEP.
- **Location Retrieval**: Finds the city name associated with the CEP.
- **Temperature Data**: Fetches current temperature in Celsius and converts it to Fahrenheit and Kelvin.
- **Error Handling**: Provides appropriate HTTP responses for different failure scenarios.
- **Dockerized**: Comes with a Dockerfile for easy containerization.
- **Cloud Deployment**: Ready for deployment on Google Cloud Run (free tier).


## Prerequisites

- **Go**: Version 1.19 or later.
- **Docker**: For containerization.
- **Google Cloud SDK**: For deployment to Google Cloud Run.
- **WeatherAPI Key**: Obtain from [WeatherAPI](https://www.weatherapi.com/).


## Installation

1. **Clone the Repository**

   ```bash
   git clone https://github.com/Twsouza/go-weather-service-api.git
   cd weather-app
   ```

2. **Set the WeatherAPI Key**

   Export your WeatherAPI key as an environment variable:

   ```bash
   export WEATHER_API_KEY=your_api_key
   ```

3. **Download Dependencies**

   ```bash
   go mod download
   ```

4. **Build the Application**

   ```bash
   go build -o weather-app
   ```

## Usage

### Running Locally

1. **Start the Application**

   ```bash
   ./weather-app
   ```

2. **Make a Request**

   Use `curl` or a web browser:

   ```bash
   curl 'http://localhost:8080/weather?cep=01001000'
   ```

## API Endpoints

- **Endpoint**: `/weather`
- **Method**: `GET`
- **Query Parameters**:
  - `cep`: The 8-digit Brazilian postal code.

### Success Response

- **HTTP Status Code**: `200 OK`
- **Response Body**:

  ```json
  {
    "temp_C": 28.5,
    "temp_F": 83.3,
    "temp_K": 301.5
  }
  ```

### Error Responses

1. **Invalid CEP Format**

   - **HTTP Status Code**: `422 Unprocessable Entity`
   - **Response Body**: `invalid zipcode`

2. **CEP Not Found**

   - **HTTP Status Code**: `404 Not Found`
   - **Response Body**: `can not find zipcode`


## Testing

### Running Tests

Execute the automated tests:

```bash
go test ./...
```


## Dockerization

### Building the Docker Image

```bash
docker build -t weather-app .
```

### Running the Docker Container

```bash
docker run -p 8080:8080 -e WEATHER_API_KEY=your_api_key weather-app
```


## Deployment on Google Cloud Run

### Prerequisites

- **Google Cloud Account**: Sign up [here](https://cloud.google.com/).
- **Google Cloud SDK**: Install from [here](https://cloud.google.com/sdk/docs/install).

### Steps

1. **Authenticate with Google Cloud**

   ```bash
   gcloud auth login
   ```

2. **Set Your Project**

   ```bash
   gcloud config set project YOUR_PROJECT_ID
   ```

3. **Build and Push the Docker Image**

   ```bash
   gcloud builds submit --tag gcr.io/YOUR_REPOSITORY/weather-app
   ```

4. **Deploy to Cloud Run**

   ```bash
   gcloud run deploy weather-app \
     --image gcr.io/YOUR_REPOSITORY/weather-app \
     --platform managed \
     --region us-central1 \
     --allow-unauthenticated \
     --set-env-vars WEATHER_API_KEY=your_api_key
   ```

5. **Access Your Application**

   After deployment, you'll receive a URL where your service is hosted.
