package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/boma026/BarberFlow-API/internal/db"
	"github.com/boma026/BarberFlow-API/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/lib/pq" // Driver do Postgres
)

func main() {
	// 1. Conexão com o banco de dados (Ajuste a senha de acordo com o seu Postgres local)
	connStr := "host=db user=postgres password=rootpassword dbname=barberflow sslmode=disable"
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Não foi possível conectar ao banco:", err)
	}
	defer conn.Close()

	// 2. Inicializa o sqlc com a conexão
	queries := db.New(conn)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// 3. Rotas com CRUD e Relacionamento
	r.Post("/customers", handlers.CreateCustomer(queries))
	r.Get("/customers", handlers.ListCustomers(queries))
	r.Get("/customers/{id}", handlers.GetCustomerNested(queries)) // O endpoint da nota máxima!
	r.Put("/customers/{id}", handlers.UpdateCustomer(queries))
	r.Delete("/customers/{id}", handlers.DeleteCustomer(queries))

	r.Post("/appointments", handlers.CreateAppointment(queries))

	println("Servidor BarberFlow rodando em http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
