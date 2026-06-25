package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/boma026/BarberFlow-API/internal/db"
	"github.com/boma026/BarberFlow-API/internal/handlers"
	"github.com/boma026/BarberFlow-API/internal/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	_ "github.com/lib/pq"
)

func main() {

	connStr := "host=db user=postgres password=rootpassword dbname=barberflow sslmode=disable"
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Não foi possível conectar ao banco:", err)
	}
	defer conn.Close()

	queries := db.New(conn)

	r := chi.NewRouter()

	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.RateLimiter())

	r.Get("/health", handlers.HealthCheck)
	r.Post("/register", handlers.RegisterUser(queries))
	r.Post("/login", handlers.LoginUser(queries))
	r.Post("/refresh", handlers.RefreshToken(queries))

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Post("/customers", handlers.CreateCustomer(queries))
		r.Get("/customers", handlers.ListCustomers(queries))
		r.Get("/customers/{id}", handlers.GetCustomerNested(queries))

		r.Post("/appointments", handlers.CreateAppointment(queries))
	})

	println("Servidor BarberFlow rodando em http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
