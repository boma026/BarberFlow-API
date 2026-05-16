package handlers

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/boma026/BarberFlow-API/internal/db"
	_ "github.com/lib/pq" // Driver do Postgres
)

// setupTestDB é uma função auxiliar para conectar ao banco do Docker durante os testes
func setupTestDB(t *testing.T) *db.Queries {
	// Importante: Note que o host é localhost, pois o teste roda no seu Windows
	// apontando para a porta 5432 exposta pelo Docker.
	connStr := "host=localhost user=postgres password=rootpassword dbname=barberflow sslmode=disable"
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Não foi possível conectar ao banco de testes: %v", err)
	}
	return db.New(conn)
}

// 1. Teste de Rota Simples (Sem Banco de Dados)
// Nota: Adicione a função HealthCheck de volta no seu handlers.go se você a tiver apagado!
func TestHealthCheck(t *testing.T) {
	req, _ := http.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	// Como o HealthCheck não acessa o banco, passamos a função diretamente
	handler := http.HandlerFunc(HealthCheck)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("HealthCheck retornou status errado: obtido %v esperado %v", status, http.StatusOK)
	}
}

// 2. Teste de Integração Real: Listar Clientes
func TestListCustomers(t *testing.T) {
	q := setupTestDB(t) // Conecta no banco real do Docker

	req, _ := http.NewRequest("GET", "/customers", nil)
	rr := httptest.NewRecorder()

	// Injetamos a conexão do banco na hora de chamar o handler
	handler := ListCustomers(q)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("ListCustomers retornou status errado: obtido %v esperado %v", status, http.StatusOK)
	}
}

// 3. Teste de Validação de Erro (Garante que a API não aceita lixo)
func TestCreateCustomer_InvalidJSON(t *testing.T) {
	q := setupTestDB(t)

	// Enviando um JSON quebrado de propósito
	jsonStr := []byte(`{"name":"Arthur" Email"faltando aspas"}`)
	req, _ := http.NewRequest("POST", "/customers", bytes.NewBuffer(jsonStr))
	rr := httptest.NewRecorder()

	handler := CreateCustomer(q)
	handler.ServeHTTP(rr, req)

	// Esperamos um erro 400 Bad Request, pois o JSON é inválido
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("CreateCustomer com JSON inválido deveria retornar 400, obteve %v", status)
	}
}
