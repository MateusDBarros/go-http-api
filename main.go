package main

import (
	"errors"
	"log"
	"net/http"
	"sync"
)

var ErrPersonNotFound = errors.New("person not found")
var ErrPersonAlreadyExists = errors.New("person already exists")

type PersonRepository interface {
	Create(p Person) error
	GetByID(id string) (Person, error)
	Update(p Person) error
	Delete(id string) error
	List() ([]Person, error)
}
type Person struct {
	ID   string
	Name string
	Age  int
	Job  string
}

type PersonStore struct {
	mu     sync.RWMutex
	people map[string]Person
}

func NewPersonStore() *PersonStore {
	return &PersonStore{
		people: make(map[string]Person),
	}
}

func (s *PersonStore) Create(p Person) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if p.ID == "" {
		return errors.New("person ID cannot be empty")
	}
	if _, exists := s.people[p.ID]; exists {
		return ErrPersonAlreadyExists
	}
	s.people[p.ID] = p
	return nil
}

func (s *PersonStore) GetByID(id string) (Person, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.people[id]
	if !ok {
		return p, ErrPersonNotFound
	}
	return p, nil
}

func (s *PersonStore) Update(p Person) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.people[p.ID]
	if !ok {
		return ErrPersonNotFound
	}
	s.people[p.ID] = p
	return nil
}

func (s *PersonStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.people[id]
	if !ok {
		return ErrPersonNotFound
	}
	delete(s.people, id)
	return nil
}
func (s *PersonStore) List() ([]Person, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ps := make([]Person, 0, len(s.people))
	for _, p := range s.people {
		ps = append(ps, p)
	}
	return ps, nil
}

func main() {
	// Initializing a thread safe store
	store := NewPersonStore()
	// Initializing the handler with store injection
	handler := NewPersonHandler(store)

	// Creating a new route
	mux := http.NewServeMux()

	// Register the routes
	mux.HandleFunc("POST /people", handler.handleCreate)
	mux.HandleFunc("GET /people/{id}", handler.handleGet)
	mux.HandleFunc("DELETE /people/{id}", handler.handleDelete)
	mux.HandleFunc("PUT /people/{id}", handler.handleUpdate)

	port := "8080"
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
