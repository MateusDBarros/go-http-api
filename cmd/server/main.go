package main

import (
	"crud/internal/api"
	"crud/internal/person"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	connStr := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping postgres: %v", err)
	}

	repo := person.NewPostgresStore(db)
	handler := api.NewPersonHandler(repo)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	port := "8080"
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
