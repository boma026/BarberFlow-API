package handlers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/boma026/BarberFlow-API/internal/db"
	_ "github.com/lib/pq"
)

// setupTestDB conecta ao banco durante os testes
func setupTestDB(t *testing.T) *db.Queries {
	connStr := "host=localhost user=postgres password=rootpassword dbname=barberflow sslmode=disable"
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Não foi possível conectar ao banco de testes: %v", err)
	}
	return db.New(conn)
}

func TestHealthCheck(t *testing.T) {
	req, _ := http.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(HealthCheck)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("HealthCheck retornou status %v, esperado %v", status, http.StatusOK)
	}
}

func TestListCustomers(t *testing.T) {
	q := setupTestDB(t)

	req, _ := http.NewRequest("GET", "/customers", nil)
	rr := httptest.NewRecorder()

	handler := ListCustomers(q)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("ListCustomers retornou status %v, esperado %v", status, http.StatusOK)
	}
}
