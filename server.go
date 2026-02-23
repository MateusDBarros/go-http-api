package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

type PersonHandler struct {
	repo PersonRepository
}

func NewPersonHandler(repo PersonRepository) *PersonHandler {
	return &PersonHandler{repo}
}

func (h *PersonHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var p Person

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.repo.Create(p); err != nil {
		if errors.Is(err, ErrPersonAlreadyExists) {
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
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	person, err := h.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, ErrPersonNotFound) {
			http.Error(w, "person not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(person)
}

func (h *PersonHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.repo.Delete(id); err != nil {
		if errors.Is(err, ErrPersonNotFound) {
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
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	var p Person

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	p.ID = id
	if err := h.repo.Update(p); err != nil {
		if errors.Is(err, ErrPersonNotFound) {
			http.Error(w, "person not found", http.StatusNotFound)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
