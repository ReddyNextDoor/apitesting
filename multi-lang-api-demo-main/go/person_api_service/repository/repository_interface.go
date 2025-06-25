package repository

import "github.com/example/person_api_service/models"

// PersonRepositoryInterface defines the operations for person data storage.
type PersonRepositoryInterface interface {
	CreatePerson(person models.PersonCreate) (*models.Person, error)
	GetPerson(id string) (*models.Person, error)
	UpdatePerson(id string, personData models.PersonUpdate) (*models.Person, error)
	DeletePerson(id string) (bool, error)
	SearchByName(firstName, lastName string) ([]models.Person, error)
	ListByCityState(city, state string) ([]models.Person, error)
}
