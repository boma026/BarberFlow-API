package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Estrutura para os dados do cliente (P1 do Backlog)
type CustomerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Estrutura para a criação de serviços (Escopo do MVP: Corte, Barba)
type ServiceRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// Estrutura para representação de saída de serviços
type Service struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// Endpoint 1: Health Check
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("BarberFlow API está online e operante!"))
}

// Endpoint 2: Listar Serviços (P1 do Backlog)
func ListServices(w http.ResponseWriter, r *http.Request) {
	// Mock
	services := []Service{
		{ID: 1, Name: "Corte de Cabelo Masculino", Price: 35.0},
		{ID: 2, Name: "Barba com Toalha Quente", Price: 25.0},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

// Endpoint 3: Criar Novo Serviço (Agora processando JSON)
func CreateService(w http.ResponseWriter, r *http.Request) {
	var req ServiceRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Erro ao processar JSON do serviço", http.StatusBadRequest)
		return
	}

	responseMessage := fmt.Sprintf("Serviço '%s' (R$ %.2f) adicionado ao catálogo! (Mock)", req.Name, req.Price)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(responseMessage))
}

// Endpoint 4: Registro de Clientes (P1 do Backlog)
func RegisterCustomer(w http.ResponseWriter, r *http.Request) {
	var req CustomerRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Erro ao processar JSON do cliente", http.StatusBadRequest)
		return
	}

	responseMessage := fmt.Sprintf("Cliente %s cadastrado com sucesso! (Mock)", req.Name)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(responseMessage))
}
