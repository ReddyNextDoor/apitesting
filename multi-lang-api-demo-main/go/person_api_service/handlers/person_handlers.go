package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/example/person_api_service/models"
	"github.com/example/person_api_service/repository"
	"github.com/gin-gonic/gin"
)

// PersonHandler holds the repository implementation.
type PersonHandler struct {
	Repo repository.PersonRepositoryInterface
}

// NewPersonHandler creates a new PersonHandler with the given repository.
// This function will decide which repository implementation to use based on ENV.
func NewPersonHandler() *PersonHandler {
	var repo repository.PersonRepositoryInterface
	backend := strings.ToLower(os.Getenv("PERSON_REPO_BACKEND"))

	log.Printf("Selected backend: %s", backend)

	if backend == "mongo" {
		log.Println("Using MongoDB repository")
		repo = repository.NewMongoPersonRepository()
	} else { // Default to SQLite
		log.Println("Using SQLite repository (default)")
		repo = repository.NewSQLitePersonRepository()
	}
	return &PersonHandler{Repo: repo}
}

// CreatePerson godoc
// @Summary Create a new person
// @Description Add a new person to the database
// @Tags persons
// @Accept  json
// @Produce  json
// @Param   person  body   models.PersonCreate  true  "Person to create"
// @Success 201 {object} models.PersonOut
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 409 {object} map[string]string "Email already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /persons/ [post]
func (h *PersonHandler) CreatePerson(c *gin.Context) {
	var personCreate models.PersonCreate
	if err := c.ShouldBindJSON(&personCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Basic validation (more can be added)
	if personCreate.FirstName == "" || personCreate.LastName == "" || personCreate.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "First name, last name, and email are required"})
		return
	}

	createdPerson, err := h.Repo.CreatePerson(personCreate)
	if err != nil {
		if err.Error() == "email already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			log.Printf("Error creating person: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create person"})
		}
		return
	}
	c.JSON(http.StatusCreated, createdPerson.ToPersonOut())
}

// GetPerson godoc
// @Summary Get a person by ID
// @Description Retrieve a person's details by their ID
// @Tags persons
// @Produce  json
// @Param   person_id  path   string  true  "Person ID"
// @Success 200 {object} models.PersonOut
// @Failure 400 {object} map[string]string "Invalid ID format"
// @Failure 404 {object} map[string]string "Person not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /persons/{person_id} [get]
func (h *PersonHandler) GetPerson(c *gin.Context) {
	id := c.Param("person_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Person ID is required"})
		return
	}

	person, err := h.Repo.GetPerson(id)
	if err != nil {
		if err.Error() == "invalid ID format" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			log.Printf("Error getting person by ID %s: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve person"})
		}
		return
	}

	if person == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
		return
	}
	c.JSON(http.StatusOK, person.ToPersonOut())
}

// UpdatePerson godoc
// @Summary Update an existing person
// @Description Modify the details of an existing person by their ID
// @Tags persons
// @Accept  json
// @Produce  json
// @Param   person_id  path   string  true  "Person ID"
// @Param   person  body   models.PersonUpdate  true  "Person data to update"
// @Success 200 {object} models.PersonOut
// @Failure 400 {object} map[string]string "Invalid input or ID format"
// @Failure 404 {object} map[string]string "Person not found"
// @Failure 409 {object} map[string]string "Email already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /persons/{person_id} [put]
func (h *PersonHandler) UpdatePerson(c *gin.Context) {
	id := c.Param("person_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Person ID is required"})
		return
	}

	var personUpdate models.PersonUpdate
	if err := c.ShouldBindJSON(&personUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	updatedPerson, err := h.Repo.UpdatePerson(id, personUpdate)
	if err != nil {
		if err.Error() == "invalid ID format" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if err.Error() == "email already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			log.Printf("Error updating person ID %s: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update person"})
		}
		return
	}

	if updatedPerson == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
		return
	}
	c.JSON(http.StatusOK, updatedPerson.ToPersonOut())
}

// DeletePerson godoc
// @Summary Delete a person by ID
// @Description Remove a person from the database by their ID
// @Tags persons
// @Produce  json
// @Param   person_id  path   string  true  "Person ID"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} map[string]string "Invalid ID format"
// @Failure 404 {object} map[string]string "Person not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /persons/{person_id} [delete]
func (h *PersonHandler) DeletePerson(c *gin.Context) {
	id := c.Param("person_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Person ID is required"})
		return
	}

	success, err := h.Repo.DeletePerson(id)
	if err != nil {
		if err.Error() == "invalid ID format" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			log.Printf("Error deleting person ID %s: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete person"})
		}
		return
	}

	if !success {
		c.JSON(http.StatusNotFound, gin.H{"error": "Person not found or already deleted"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// SearchPersons godoc
// @Summary Search for persons by name
// @Description Search for persons by first name and/or last name (case-insensitive, partial match)
// @Tags persons
// @Produce  json
// @Param   first_name  query  string  false  "First name to search for"
// @Param   last_name   query  string  false  "Last name to search for"
// @Success 200 {array}  models.PersonOut
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /persons/search [get]
func (h *PersonHandler) SearchPersons(c *gin.Context) {
	firstName := c.Query("first_name")
	lastName := c.Query("last_name")

	// Optional: Add validation if at least one parameter is required
	// if firstName == "" && lastName == "" {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "At least one search parameter (first_name or last_name) is required"})
	// 	return
	// }

	persons, err := h.Repo.SearchByName(firstName, lastName)
	if err != nil {
		log.Printf("Error searching persons (FN: %s, LN: %s): %v", firstName, lastName, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search persons"})
		return
	}

	personsOut := make([]models.PersonOut, len(persons))
	for i, p := range persons {
		personsOut[i] = p.ToPersonOut()
	}
	c.JSON(http.StatusOK, personsOut)
}

// PersonsByCityState godoc
// @Summary List persons by city and state
// @Description Retrieve a list of persons residing in a specific city and state (case-insensitive, exact match)
// @Tags persons
// @Produce  json
// @Param   city   query  string  true  "City name"
// @Param   state  query  string  true  "State name"
// @Success 200 {array}  models.PersonOut
// @Failure 400 {object} map[string]string "City and state are required"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /persons/by_city_state [get]
func (h *PersonHandler) PersonsByCityState(c *gin.Context) {
	city := c.Query("city")
	state := c.Query("state")

	if city == "" || state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "City and state query parameters are required"})
		return
	}

	persons, err := h.Repo.ListByCityState(city, state)
	if err != nil {
		if err.Error() == "city and state must be provided" { // This check might be redundant if we validate above
			c.JSON(http.StatusBadRequest, gin.H{"error": "City and state query parameters are required"})
		} else {
			log.Printf("Error listing persons by city '%s' and state '%s': %v", city, state, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list persons by city and state"})
		}
		return
	}

	personsOut := make([]models.PersonOut, len(persons))
	for i, p := range persons {
		personsOut[i] = p.ToPersonOut()
	}
	c.JSON(http.StatusOK, personsOut)
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Returns the current backend being used by the service
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	backend := strings.ToLower(os.Getenv("PERSON_REPO_BACKEND"))
	if backend == "" {
		backend = "sqlite" // Default if not set
	}

	mongoDbName := os.Getenv("MONGO_DB")
	if mongoDbName == "" {
		mongoDbName = "person_db"
	}

	sqliteDbPath := os.Getenv("SQLITE_DB_PATH")
	if sqliteDbPath == "" {
		sqliteDbPath = "db/persons.db"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":         "ok",
		"backend":        backend,
		"mongo_db_name":  mongoDbName,
		"sqlite_db_path": sqliteDbPath,
	})
}
