package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/sirupsen/logrus"
)

// ZipCodeService defines the interface for fetching location by ZIP code.
type ZipCodeService interface {
	GetLocationByZipCode(zipCode string) (string, error)
	IsValidCEP(cep string) bool
}

// ViaCEPService implements ZipCodeService using the ViaCEP API.
type ViaCEPService struct {
	BaseURL string
}

func NewViaCEPService(baseURL string) *ViaCEPService {
	return &ViaCEPService{
		BaseURL: baseURL,
	}
}

// viaCEPResponse represents the response from the ViaCEP API.
type viaCEPResponse struct {
	Localidade string `json:"localidade"`
}

// GetLocationByZipCode fetches the city name for a given ZIP code.
func (v *ViaCEPService) GetLocationByZipCode(zipCode string) (string, error) {
	apiURL, err := url.Parse(v.BaseURL)
	if err != nil {
		return "", err
	}
	apiURL.Path = fmt.Sprintf("/ws/%s/json", zipCode)

	logrus.WithFields(logrus.Fields{
		"apiURL": apiURL.String(),
	}).Info("Fetching location by ZIP code", zipCode)

	resp, err := http.Get(apiURL.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("can not find zipcode")
	}

	var data viaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	logrus.WithFields(logrus.Fields{
		"city": data.Localidade,
	}).Info("City found for ZIP code", zipCode)

	return data.Localidade, nil
}

// IsValidCEP validates if the CEP is a valid 8-digit number.
func (v *ViaCEPService) IsValidCEP(cep string) bool {
	matched, err := regexp.MatchString(`^\d{8}$`, cep)
	if err != nil {
		logrus.Error("Error while validating CEP: ", err)

		return false
	}

	return matched
}
