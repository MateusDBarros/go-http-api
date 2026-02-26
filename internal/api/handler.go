package api

import (
	"crud/internal/person"
	"encoding/json"
	"errors"
	"net/http"
)

type PersonHandler struct {
	repo person.Repository
}

func NewPersonHandler(repo person.Repository) *PersonHandler {
	return &PersonHandler{repo: repo}
}

func (h *PersonHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /people", h.handleCreate)
	mux.HandleFunc("GET /people/{id}", h.handleGet)
	mux.HandleFunc("PUT /people/{id}", h.handleUpdate)
	mux.HandleFunc("DELETE /people/{id}", h.handleDelete)
}

func (h *PersonHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var p person.Person
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(p); err != nil {
		if errors.Is(err, person.ErrPersonAlreadyExists) {
			http.Error(w, "person already exists", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *PersonHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "invalid person id", http.StatusBadRequest)
		return
	}

	p, err := h.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, person.ErrPersonNotFound) {
			http.Error(w, "person not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *PersonHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "invalid person id", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		if errors.Is(err, person.ErrPersonNotFound) {
			http.Error(w, "person not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *PersonHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "invalid person id", http.StatusBadRequest)
		return
	}

	var p person.Person
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	p.ID = id

	if err := h.repo.Update(p); err != nil {
		if errors.Is(err, person.ErrPersonNotFound) {
			http.Error(w, "person not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
