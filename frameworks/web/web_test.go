package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/riomhaire/lightauth2/entities"
	"github.com/riomhaire/lightauth2/frameworks"
	"github.com/riomhaire/lightauth2/interfaces"
	"github.com/riomhaire/lightauth2/usecases"
)

func createTestRegistry() usecases.Registry {
	logger := frameworks.ConsoleLogger{}

	database := frameworks.NewStringDatabaseInteractor()
	database.Create(entities.User{"test", "939c1f673b7f5f5c991b4d4160642e72880e783ba4d7b04da260392f855214a6", true, []string{"user"}, "claim1", "claim2"})
	database.Create(entities.User{"admin", "50b911deac5df04e0a79ef18b04b29b245b8f576dcb7e5cca5937eb2083438ba", true, []string{"admin"}, "claim1", "claim2"})

	configuration := usecases.Configuration{}
	configuration.SigningSecret = "secret"
	configuration.TokenTimeout = 3600

	registry := usecases.Registry{}
	registry.Configuration = configuration
	registry.Logger = logger
	registry.StorageInteractor = database
	registry.AuthenticateInteractor = interfaces.DefaultAuthenticateInteractor{&registry}
	registry.TokenInteractor = interfaces.DefaultTokenInteractor{&registry}

	return registry
}

func TestAuthenticateHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	body := []byte("{\"username\":\"test\",\"password\":\"secret\"}")
	req, err := http.NewRequest("POST", "/authenticate", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	registry := createTestRegistry()
	restAPI := NewRestAPI(&registry)
	handler := http.HandlerFunc(restAPI.HandleAuthenticate)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	authenticateResponse := AuthenticateResponse{}
	err = json.NewDecoder(rr.Body).Decode(&authenticateResponse)
	if err != nil {
		t.Fatal(err)
	}
	// Check we have a popoulated token field
	if authenticateResponse.Token == "" || len(authenticateResponse.Token) < 10 {
		t.Fatal(errors.New("token missing or invalid"))
	}
}

func TestAuthenticateHandlerBadCredentials(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	body := []byte("{\"username\":\"test\",\"password\":\"test\"}")
	req, err := http.NewRequest("POST", "/authenticate", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	registry := createTestRegistry()
	restAPI := NewRestAPI(&registry)
	handler := http.HandlerFunc(restAPI.HandleAuthenticate)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

}

func TestValidateHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	body := []byte("{\"username\":\"test\",\"password\":\"secret\"}")
	req, err := http.NewRequest("POST", "/authenticate", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	registry := createTestRegistry()
	restAPI := NewRestAPI(&registry)
	handler := http.HandlerFunc(restAPI.HandleAuthenticate)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	authenticateResponse := AuthenticateResponse{}
	err = json.NewDecoder(rr.Body).Decode(&authenticateResponse)
	if err != nil {
		t.Fatal(err)
	}
	// Check we have a popoulated token field
	if authenticateResponse.Token == "" || len(authenticateResponse.Token) < 10 {
		t.Fatal(errors.New("token missing or invalid"))
	}

	// OK Validate
	authorization := fmt.Sprintf("Bearer %s", authenticateResponse.Token)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(restAPI.HandleValidate)
	req, _ = http.NewRequest("GET", "/api/v2/session", nil)
	req.Header.Set("Authorization", authorization)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestTokenDecodeHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	body := []byte("{\"username\":\"test\",\"password\":\"secret\"}")
	req, err := http.NewRequest("POST", "/authenticate", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	registry := createTestRegistry()
	restAPI := NewRestAPI(&registry)
	handler := http.HandlerFunc(restAPI.HandleAuthenticate)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	authenticateResponse := AuthenticateResponse{}
	err = json.NewDecoder(rr.Body).Decode(&authenticateResponse)
	if err != nil {
		t.Fatal(err)
	}
	// Check we have a popoulated token field
	if authenticateResponse.Token == "" || len(authenticateResponse.Token) < 10 {
		t.Fatal(errors.New("token missing or invalid"))
	}

	// OK Decode token
	authorization := fmt.Sprintf("Bearer %s", authenticateResponse.Token)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(restAPI.HandleTokenDecode)
	req, _ = http.NewRequest("GET", "/api/v2/session/decoder", nil)
	req.Header.Set("Authorization", authorization)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}

func TestTokenDecodeHandlerNoBearer(t *testing.T) {
	// OK Decode token no auth
	registry := createTestRegistry()
	restAPI := NewRestAPI(&registry)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(restAPI.HandleTokenDecode)
	req, _ := http.NewRequest("GET", "/api/v2/session/decoder", nil)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}

func TestValidateNoAuthorizationHeaderHandler(t *testing.T) {
	// OK Validate
	rr := httptest.NewRecorder()
	registry := createTestRegistry()
	restAPI := NewRestAPI(&registry)
	handler := http.HandlerFunc(restAPI.HandleValidate)
	req, _ := http.NewRequest("GET", "/api/v2/session", nil)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestValidateBadAuthorizationHeaderHandler(t *testing.T) {
	// OK Validate
	rr := httptest.NewRecorder()
	registry := createTestRegistry()
	restAPI := NewRestAPI(&registry)
	handler := http.HandlerFunc(restAPI.HandleValidate)
	req, _ := http.NewRequest("GET", "/api/v2/session", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer BAD"))
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestHealthHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/api/v2/authentication/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	registry := createTestRegistry()
	restAPI := NewRestAPI(&registry)
	handler := http.HandlerFunc(restAPI.HandleHealth)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}

func TestStatisticsHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/api/v2/authentication/statistics", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	registry := createTestRegistry()
	restAPI := NewRestAPI(&registry)
	handler := http.HandlerFunc(restAPI.HandleStatistics)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check for presence of common fields pid and uptime
	responseMap := make(map[string]interface{})
	err = json.NewDecoder(rr.Body).Decode(&responseMap)
	if err != nil {
		t.Fatal(err)
	}

	for _, val := range []string{"pid", "uptime"} {
		if _, ok := responseMap[val]; ok {
			//prsent
		} else {
			// missing - fail
			t.Fatal(errors.New("Expected parameter missing"))
		}

	}

}

func TestPrometheusStatisticsHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/api/v2/authentication/statistics", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Set Accept type to text
	req.Header.Set("Accept", "text/plain")

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	registry := createTestRegistry()
	restAPI := NewRestAPI(&registry)
	handler := http.HandlerFunc(restAPI.HandleStatistics)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Read response and check has content
	body := string(rr.Body.Bytes())
	if len(body) == 0 {
		t.Error("Expected String back but its empty")
	}

	// Check lightauth2_response_total_count present
	if !strings.Contains(body, "lightauth2_response_total_count") {
		t.Errorf("Expected 'lightauth2_response_total_count' but its not there. I got %v\n", body)
	}

}

func TestHostAppender(t *testing.T) {
	// Check that extra headers are added
	rr := httptest.NewRecorder()
	api := RestAPI{}
	api.AddWorkerHeader(rr, nil, nil)
	if len(rr.Header().Get("X-WORKER")) == 0 {
		t.Fail()
	}
	// Check that non-existent header not there
	if len(rr.Header().Get("X-WORKER2")) != 0 {
		t.Fail()
	}
}

func TestVersionAppender(t *testing.T) {
	// Check that extra headers are added
	rr := httptest.NewRecorder()
	registry := createTestRegistry()
	api := NewRestAPI(&registry)
	api.AddWorkerVersion(rr, nil, nil)
	if len(rr.Header().Get("X-WORKER-VERSION")) == 0 {
		t.Fail()
	}
	// Check that non-existent header not there
	if len(rr.Header().Get("X-WORKER2-VERSION")) != 0 {
		t.Fail()
	}
}
