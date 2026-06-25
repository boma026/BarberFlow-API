package handlers

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/boma026/BarberFlow-API/internal/auth"
	"github.com/boma026/BarberFlow-API/internal/db"
	"github.com/boma026/BarberFlow-API/internal/middleware"
	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *db.Queries {
	connStr := "host=localhost user=postgres password=rootpassword dbname=barberflow sslmode=disable"
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Não foi possível conectar ao banco de testes: %v", err)
	}
	return db.New(conn)
}

// 1. Teste de Rota Pública Simples
func TestHealthCheck(t *testing.T) {
	req, _ := http.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheck)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Esperado 200, obtido %v", status)
	}
}

// 2. Teste da Proteção OWASP (Security Headers)
func TestSecurityHeaders_Injected(t *testing.T) {
	req, _ := http.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	// Passa a rota por dentro do middleware de segurança
	handler := middleware.SecurityHeaders(http.HandlerFunc(HealthCheck))
	handler.ServeHTTP(rr, req)

	if rr.Header().Get("X-Frame-Options") != "DENY" {
		t.Errorf("Cabeçalho de segurança OWASP não foi injetado")
	}
}

// 3. Teste de Bloqueio sem Token (Auth Middleware)
func TestAccessProtectedWithoutToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "/customers", nil)
	rr := httptest.NewRecorder()

	q := setupTestDB(t)
	handler := middleware.AuthMiddleware(ListCustomers(q))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Esperado 401 Unauthorized para acesso sem token, obtido %v", status)
	}
}

// 4. Teste de Bloqueio com Token Falso
func TestAccessProtectedWithInvalidToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "/customers", nil)
	req.Header.Set("Authorization", "Bearer token_completamente_falso_e_invalido")
	rr := httptest.NewRecorder()

	q := setupTestDB(t)
	handler := middleware.AuthMiddleware(ListCustomers(q))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Esperado 401 Unauthorized para token inválido, obtido %v", status)
	}
}

// 5. Teste do Gerador de Tokens JWT
func TestAuthTokenGeneration(t *testing.T) {
	email := "teste_unitario@ufrn.br"
	access, refresh, err := auth.GenerateTokens(email)

	if err != nil || access == "" || refresh == "" {
		t.Errorf("Falha do motor ao gerar os tokens JWT")
	}
}

// 6. Teste de Login com Corpo Vazio
func TestLoginUser_InvalidBody(t *testing.T) {
	q := setupTestDB(t)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(`{}`)))
	rr := httptest.NewRecorder()

	handler := LoginUser(q)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Esperado 401 Unauthorized para credenciais vazias, obtido %v", status)
	}
}

func TestRegisterUser_InvalidJSON(t *testing.T) {
	q := setupTestDB(t)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(`{"email": "incompleto"`)))
	rr := httptest.NewRecorder()

	handler := RegisterUser(q)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Esperado 400 Bad Request para JSON quebrado, obtido %v", status)
	}
}

func TestRefreshToken_Invalid(t *testing.T) {
	q := setupTestDB(t)
	req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer([]byte(`{"refresh_token": "falso"}`)))
	rr := httptest.NewRecorder()

	handler := RefreshToken(q)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Esperado 401 Unauthorized para tentar renovar com token falso, obtido %v", status)
	}
}

func TestCreateCustomer_InvalidJSON(t *testing.T) {
	q := setupTestDB(t)
	jsonStr := []byte(`{"name":"Arthur" Email"faltando aspas"}`)
	req, _ := http.NewRequest("POST", "/customers", bytes.NewBuffer(jsonStr))
	rr := httptest.NewRecorder()

	handler := CreateCustomer(q)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Esperado 400 Bad Request, obtido %v", status)
	}
}

func TestGetCustomerNested_InvalidID(t *testing.T) {
	q := setupTestDB(t)

	req, _ := http.NewRequest("GET", "/customers/letras", nil)
	rr := httptest.NewRecorder()

	handler := GetCustomerNested(q)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Esperado 400 Bad Request ao tentar buscar ID com formato errado, obtido %v", status)
	}
}
