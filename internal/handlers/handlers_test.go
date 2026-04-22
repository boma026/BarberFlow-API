package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHealthCheck verifica se a API responde 200 OK
func TestHealthCheck(t *testing.T) {
	req, _ := http.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheck)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler retornou status errado: obtido %v esperado %v", status, http.StatusOK)
	}
}

// TestListServices verifica se a listagem retorna 200 OK
func TestListServices(t *testing.T) {
	req, _ := http.NewRequest("GET", "/services", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ListServices)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler retornou status errado: obtido %v esperado %v", status, http.StatusOK)
	}
}

// TestRegisterCustomer verifica se o cadastro processa o JSON corretamente
func TestRegisterCustomer(t *testing.T) {
	jsonStr := []byte(`{"name":"Arthur Boma","email":"arthur@ufrn.br"}`)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonStr))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RegisterCustomer)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler retornou status errado: obtido %v esperado %v", status, http.StatusCreated)
	}
}
