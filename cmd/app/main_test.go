package main

import (
	"bytes"
	"csye6225-mainproject/routes"
	"csye6225-mainproject/services"
	"errors"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockHealthzStore struct {
	db            *gorm.DB
	isDBInService bool
}

func (m *MockHealthzStore) OpenDBConnection(dialector gorm.Dialector, config *gorm.Config) error {

	if m.isDBInService {
		m.db = &gorm.DB{}
		return nil
	} else {
		m.db = nil
		NotAbleToConnectToDBError := errors.New("not able to connect to database")
		return NotAbleToConnectToDBError
	}

}

func (m *MockHealthzStore) GetDBConnection() *gorm.DB {
	return m.db
}

func (m *MockHealthzStore) CloseDBConnection() error {
	return nil
}

func (m *MockHealthzStore) Ping() (bool, error) {
	if m.isDBInService {
		return true, nil
	} else {
		m.db = nil
		NotAbleToConnectToDBError := errors.New("not able to connect to database")
		return false, NotAbleToConnectToDBError
	}
}

func TestMainRouter(t *testing.T) {

	t.Run("Calling /healthz endpoint should return 200 status OK when database is not in service", func(t *testing.T) {

		serviceProvider := &services.ServiceProvider{
			MyHealthzStore: &MockHealthzStore{isDBInService: true},
		}
		router := routes.SetupRouter(serviceProvider)

		request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assertErrorCode(t, response.Code, http.StatusOK)

	})

	t.Run("Calling /healthz endpoint should return 503 when database is not in service", func(t *testing.T) {

		serviceProvider := &services.ServiceProvider{
			MyHealthzStore: &MockHealthzStore{isDBInService: false},
		}
		router := routes.SetupRouter(serviceProvider)

		request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assertErrorCode(t, response.Code, http.StatusServiceUnavailable)

	})

	t.Run("Calling /healthz endpoint with every method except GET should return status 405 not allowed", func(t *testing.T) {
		serviceProvider := &services.ServiceProvider{
			MyHealthzStore: &MockHealthzStore{isDBInService: false},
		}
		router := routes.SetupRouter(serviceProvider)

		httpMethods := []string{
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodConnect,
			http.MethodHead,
			http.MethodOptions,
			http.MethodTrace,
		}

		for _, httpMethod := range httpMethods {
			request := httptest.NewRequest(httpMethod, "/healthz", nil)
			response := httptest.NewRecorder()

			router.ServeHTTP(response, request)

			assertErrorCode(t, response.Code, http.StatusMethodNotAllowed)
		}

	})

	t.Run("Accessing invalid URLs should give 404 not found", func(t *testing.T) {

		serviceProvider := &services.ServiceProvider{
			MyHealthzStore: &MockHealthzStore{isDBInService: false},
		}
		router := routes.SetupRouter(serviceProvider)

		someInvalidURLs := []string{
			"/",
			"/hello",
			"/hello/world",
			"/main",
			"/yes?name=Raja&pass=Pass",
			"/healthz/db",
		}

		for _, invalidUrl := range someInvalidURLs {
			request := httptest.NewRequest(http.MethodGet, invalidUrl, nil)
			response := httptest.NewRecorder()

			router.ServeHTTP(response, request)

			assertErrorCode(t, response.Code, http.StatusNotFound)
		}
	})

	t.Run("Check response should from any request should not have any payload", func(t *testing.T) {
		serviceProvider := &services.ServiceProvider{
			MyHealthzStore: &MockHealthzStore{isDBInService: false},
		}
		router := routes.SetupRouter(serviceProvider)

		someURLs := []string{
			"/healthz",
			"/myname",
			"/testURL?q=123",
			"/healthz/db",
		}

		for _, someUrl := range someURLs {
			request := httptest.NewRequest(http.MethodGet, someUrl, nil)
			response := httptest.NewRecorder()

			router.ServeHTTP(response, request)

			got := response.Body.String()
			want := ""

			assertString(t, got, want)
		}

	})

	t.Run("Requests with payload to valid URLs should return 400", func(t *testing.T) {

		serviceProvider := &services.ServiceProvider{
			MyHealthzStore: &MockHealthzStore{isDBInService: false},
		}
		router := routes.SetupRouter(serviceProvider)

		validURL := "/healthz"

		requestBody := &bytes.Buffer{}

		requestBody.WriteString("Some payload")
		request := httptest.NewRequest(http.MethodGet, validURL, requestBody)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assertErrorCode(t, response.Code, http.StatusBadRequest)

	})

	t.Run("Requests with valid URLs containing parameters should return 400", func(t *testing.T) {

		serviceProvider := &services.ServiceProvider{
			MyHealthzStore: &MockHealthzStore{isDBInService: false},
		}
		router := routes.SetupRouter(serviceProvider)

		someURLs := []string{
			"/healthz?q1=v1&q2=v2",
			"/healthz?hello=world&hey=man",
			"/healthz?q=123",
		}

		for _, someUrl := range someURLs {
			request := httptest.NewRequest(http.MethodGet, someUrl, nil)
			response := httptest.NewRecorder()

			router.ServeHTTP(response, request)

			assertErrorCode(t, response.Code, http.StatusBadRequest)
		}
	})

	t.Run("Requests with empty payload to valid URLs should not return 400", func(t *testing.T) {

		serviceProvider := &services.ServiceProvider{
			MyHealthzStore: &MockHealthzStore{isDBInService: false},
		}
		router := routes.SetupRouter(serviceProvider)

		validURL := "/healthz"

		request := httptest.NewRequest(http.MethodGet, validURL, nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assertNotErrorCode(t, response.Code, http.StatusBadRequest)

	})

	t.Run("Response should not be cached by setting Cache-Control to no-cache header for all URLs", func(t *testing.T) {

		serviceProvider := &services.ServiceProvider{
			MyHealthzStore: &MockHealthzStore{isDBInService: false},
		}
		router := routes.SetupRouter(serviceProvider)

		someURLs := []string{
			"/healthz",
			"/myname",
			"/testURL?q=123",
			"/healthz/db",
		}

		want := "no-cache"
		for _, someUrl := range someURLs {
			request := httptest.NewRequest(http.MethodGet, someUrl, nil)
			response := httptest.NewRecorder()

			router.ServeHTTP(response, request)

			got := response.Header().Get("Cache-Control")
			assertString(t, got, want)
		}
	})
}

func assertString(t testing.TB, got string, want string) {
	t.Helper()
	if got != want {
		t.Errorf("want %v, but got %v", want, got)
	}
}

func assertErrorCode(t testing.TB, got int, want int) {
	t.Helper()
	if got != want {
		t.Errorf("want %d, but got %d", want, got)
	}
}

func assertNotErrorCode(t testing.TB, got int, shouldNotBe int) {
	t.Helper()
	if got == shouldNotBe {
		t.Errorf("should not be %d, but got %d", shouldNotBe, got)
	}
}
