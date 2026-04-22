package main

import (
	"net/http"

	"github.com/boma026/BarberFlow-API/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", handlers.HealthCheck)
	r.Get("/services", handlers.ListServices)
	r.Post("/services", handlers.CreateService)
	r.Post("/register", handlers.RegisterCustomer)

	println("Servidor iniciado em http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
