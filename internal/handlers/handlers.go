package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/boma026/BarberFlow-API/internal/db"
	"github.com/go-chi/chi/v5"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("BarberFlow API está online e operante!"))
}

func CreateCustomer(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		customer, err := q.CreateCustomer(r.Context(), db.CreateCustomerParams{
			Name:  req.Name,
			Email: req.Email,
		})
		if err != nil {
			http.Error(w, "Erro ao salvar no banco", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(customer)
	}
}

func ListCustomers(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		customers, err := q.ListCustomers(r.Context())
		if err != nil {
			http.Error(w, "Erro ao buscar clientes", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(customers)
	}
}

func GetCustomerNested(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		customerData, err := q.GetCustomerWithAppointmentsNested(r.Context(), int32(id))
		if err != nil {
			http.Error(w, "Cliente não encontrado", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(customerData)
	}
}

func UpdateCustomer(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, _ := strconv.Atoi(idStr)

		var req struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		err := q.UpdateCustomer(r.Context(), db.UpdateCustomerParams{
			ID:    int32(id),
			Name:  req.Name,
			Email: req.Email,
		})
		if err != nil {
			http.Error(w, "Erro ao atualizar", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func DeleteCustomer(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, _ := strconv.Atoi(idStr)

		err := q.DeleteCustomer(r.Context(), int32(id))
		if err != nil {
			http.Error(w, "Erro ao deletar", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

/*
	{
	  "customer_id": 1,
	  "appointment_date": "2026-05-20T14:30:00Z",
	  "service_name": "Corte de Cabelo e Barba"
	}
*/
func CreateAppointment(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			CustomerID      int32  `json:"customer_id"`
			AppointmentDate string `json:"appointment_date"` // Ex: "2026-05-20T14:30:00Z"
			ServiceName     string `json:"service_name"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		parsedTime, _ := time.Parse(time.RFC3339, req.AppointmentDate)

		appointment, err := q.CreateAppointment(r.Context(), db.CreateAppointmentParams{
			CustomerID:      req.CustomerID,
			AppointmentDate: parsedTime,
			ServiceName:     req.ServiceName,
		})
		if err != nil {
			http.Error(w, "Erro ao criar agendamento", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(appointment)
	}
}
