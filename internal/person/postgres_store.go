package person

import (
	"database/sql"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (s *PostgresStore) Create(p Person) error {
	query := "INSERT INTO people (id, name, age, job) VALUES ($1, $2, $3, $4)"
	_, err := s.db.Exec(query, p.ID, p.Name, p.Age, p.Job)
	return err
}

func (s *PostgresStore) GetByID(id string) (Person, error) {
	query := "SELECT id, name, age, job FROM people WHERE id = $1"
	var p Person

	err := s.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Age, &p.Job)
	if err != nil {
		if err == sql.ErrNoRows {
			return p, ErrPersonNotFound
		}
		return p, err
	}
	return p, nil
}

func (s *PostgresStore) Update(p Person) error {
	query := "UPDATE people SET name = $1, age = $2, job = $3 WHERE id = $4"
	_, err := s.db.Exec(query, p.Name, p.Age, p.Job, p.ID)
	return err
}

func (s *PostgresStore) Delete(id string) error {
	query := "DELETE FROM people WHERE id = $1"
	_, err := s.db.Exec(query, id)
	return err
}

func (s *PostgresStore) List() ([]Person, error) {
	query := "SELECT id, name, age, job FROM people"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var persons []Person
	for rows.Next() {
		var p Person
		if err := rows.Scan(&p.ID, &p.Name, &p.Age, &p.Job); err != nil {
			return nil, err
		}
		persons = append(persons, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return persons, nil
}
