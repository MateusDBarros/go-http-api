package person

import (
	"errors"
	"sync"
)

type MemoryStore struct {
	mu     sync.RWMutex
	people map[string]Person
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		people: make(map[string]Person),
	}
}

func (s *MemoryStore) Create(p Person) error {
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

func (s *MemoryStore) GetByID(id string) (Person, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.people[id]
	if !ok {
		return p, ErrPersonNotFound
	}
	return p, nil
}

func (s *MemoryStore) Update(p Person) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.people[p.ID]; !ok {
		return ErrPersonNotFound
	}
	s.people[p.ID] = p
	return nil
}

func (s *MemoryStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.people[id]; !ok {
		return ErrPersonNotFound
	}
	delete(s.people, id)
	return nil
}

func (s *MemoryStore) List() ([]Person, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ps := make([]Person, 0, len(s.people))
	for _, p := range s.people {
		ps = append(ps, p)
	}
	return ps, nil
}
