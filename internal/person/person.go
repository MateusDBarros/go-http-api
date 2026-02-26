package person

import "errors"

var (
	ErrPersonNotFound      = errors.New("person not found")
	ErrPersonAlreadyExists = errors.New("person already exists")
)

type Person struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
	Job  string `json:"job"`
}

type Repository interface {
	Create(p Person) error
	GetByID(id string) (Person, error)
	Update(p Person) error
	Delete(id string) error
	List() ([]Person, error)
}
