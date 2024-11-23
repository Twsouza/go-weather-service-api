package services

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestIsValidCEP tests the IsValidCEP method.
func TestIsValidCEP(t *testing.T) {
	service := &ViaCEPService{}
	tests := []struct {
		cep     string
		isValid bool
	}{
		{"12345678", true},
		{"1234567", false},
		{"abcdefgh", false},
		{"123456789", false},
	}

	for _, test := range tests {
		result := service.IsValidCEP(test.cep)
		if result != test.isValid {
			t.Errorf("IsValidCEP(%s) = %v; want %v", test.cep, result, test.isValid)
		}
	}
}

// TestGetLocationByZipCode tests the GetLocationByZipCode method.
func TestGetLocationByZipCode(t *testing.T) {
	// Mock the HTTP response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"localidade": "Test City"}`))
	}))
	defer ts.Close()

	service := NewViaCEPService(ts.URL)

	location, err := service.GetLocationByZipCode("12345678")
	if err != nil {
		t.Fatalf("GetLocationByZipCode returned error: %v", err)
	}

	if location != "Test City" {
		t.Errorf("GetLocationByZipCode = %s; want %s", location, "Test City")
	}
}
