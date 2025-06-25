package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/example/person_api_service/database"
	"github.com/example/person_api_service/models"
	_ "github.com/mattn/go-sqlite3" // SQLite driver, ensure it's loaded
)

// SQLitePersonRepository implements PersonRepositoryInterface for SQLite.
type SQLitePersonRepository struct {
	DB *sql.DB
}

// NewSQLitePersonRepository creates a new instance of SQLitePersonRepository.
func NewSQLitePersonRepository() *SQLitePersonRepository {
	if database.SqliteDB == nil {
		database.InitSQLite() // Ensure DB is initialized
	}
	return &SQLitePersonRepository{DB: database.SqliteDB}
}

// CreatePerson adds a new person to the SQLite database.
func (r *SQLitePersonRepository) CreatePerson(personData models.PersonCreate) (*models.Person, error) {
	query := `INSERT INTO persons (first_name, last_name, email, phone, address, city, state, zip_code)
	            VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	stmt, err := r.DB.Prepare(query)
	if err != nil {
		log.Printf("Error preparing statement for CreatePerson: %v", err)
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(personData.FirstName, personData.LastName, personData.Email, personData.Phone, personData.Address, personData.City, personData.State, personData.ZipCode)
	if err != nil {
		log.Printf("Error executing statement for CreatePerson: %v", err)
		// Check for unique constraint violation for email
		if strings.Contains(err.Error(), "UNIQUE constraint failed: persons.email") {
			return nil, errors.New("email already exists")
		}
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID for CreatePerson: %v", err)
		return nil, err
	}

	createdPerson := &models.Person{
		SqliteID:  id,
		FirstName: personData.FirstName,
		LastName:  personData.LastName,
		Email:     personData.Email,
		Phone:     personData.Phone,
		Address:   personData.Address,
		City:      personData.City,
		State:     personData.State,
		ZipCode:   personData.ZipCode,
	}
	return createdPerson, nil
}

// GetPerson retrieves a person by their ID from SQLite.
func (r *SQLitePersonRepository) GetPerson(idStr string) (*models.Person, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("Invalid ID format: %s", idStr)
		return nil, errors.New("invalid ID format")
	}

	query := `SELECT id, first_name, last_name, email, phone, address, city, state, zip_code
	            FROM persons WHERE id = ?`
	row := r.DB.QueryRow(query, id)

	var p models.Person
	err = row.Scan(&p.SqliteID, &p.FirstName, &p.LastName, &p.Email, &p.Phone, &p.Address, &p.City, &p.State, &p.ZipCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Person not found
		}
		log.Printf("Error scanning person row: %v", err)
		return nil, err
	}
	return &p, nil
}

// UpdatePerson modifies an existing person's details in SQLite.
func (r *SQLitePersonRepository) UpdatePerson(idStr string, personData models.PersonUpdate) (*models.Person, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("Invalid ID format for update: %s", idStr)
		return nil, errors.New("invalid ID format")
	}

	// Check if person exists
	existingPerson, err := r.GetPerson(idStr)
	if err != nil {
		return nil, err // Error during fetch
	}
	if existingPerson == nil {
		return nil, nil // Person not found
	}

	query := "UPDATE persons SET "
	args := []interface{}{}
	updates := []string{}

	if personData.FirstName != "" {
		updates = append(updates, "first_name = ?")
		args = append(args, personData.FirstName)
		existingPerson.FirstName = personData.FirstName
	}
	if personData.LastName != "" {
		updates = append(updates, "last_name = ?")
		args = append(args, personData.LastName)
		existingPerson.LastName = personData.LastName
	}
	if personData.Email != "" {
		updates = append(updates, "email = ?")
		args = append(args, personData.Email)
		existingPerson.Email = personData.Email
	}
	if personData.Phone != "" {
		updates = append(updates, "phone = ?")
		args = append(args, personData.Phone)
		existingPerson.Phone = personData.Phone
	}
	if personData.Address != "" {
		updates = append(updates, "address = ?")
		args = append(args, personData.Address)
		existingPerson.Address = personData.Address
	}
	if personData.City != "" {
		updates = append(updates, "city = ?")
		args = append(args, personData.City)
		existingPerson.City = personData.City
	}
	if personData.State != "" {
		updates = append(updates, "state = ?")
		args = append(args, personData.State)
		existingPerson.State = personData.State
	}
	if personData.ZipCode != "" {
		updates = append(updates, "zip_code = ?")
		args = append(args, personData.ZipCode)
		existingPerson.ZipCode = personData.ZipCode
	}

	if len(updates) == 0 {
		return existingPerson, nil // No fields to update
	}

	query += strings.Join(updates, ", ") + " WHERE id = ?"
	args = append(args, id)

	stmt, err := r.DB.Prepare(query)
	if err != nil {
		log.Printf("Error preparing statement for UpdatePerson: %v", err)
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(args...)
	if err != nil {
		log.Printf("Error executing update for UpdatePerson: %v", err)
		if strings.Contains(err.Error(), "UNIQUE constraint failed: persons.email") {
			return nil, errors.New("email already exists")
		}
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected for UpdatePerson: %v", err)
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, nil // Should not happen if GetPerson check passed, but good to have
	}

	return existingPerson, nil
}

// DeletePerson removes a person from the SQLite database by ID.
func (r *SQLitePersonRepository) DeletePerson(idStr string) (bool, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("Invalid ID format for delete: %s", idStr)
		return false, errors.New("invalid ID format")
	}

	query := "DELETE FROM persons WHERE id = ?"
	stmt, err := r.DB.Prepare(query)
	if err != nil {
		log.Printf("Error preparing statement for DeletePerson: %v", err)
		return false, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		log.Printf("Error executing delete for DeletePerson: %v", err)
		return false, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected for DeletePerson: %v", err)
		return false, err
	}

	return rowsAffected > 0, nil
}

// SearchByName finds persons by first and/or last name in SQLite.
func (r *SQLitePersonRepository) SearchByName(firstName, lastName string) ([]models.Person, error) {
	var query strings.Builder
	query.WriteString("SELECT id, first_name, last_name, email, phone, address, city, state, zip_code FROM persons WHERE 1=1")
	args := []interface{}{}

	if firstName != "" {
		query.WriteString(" AND lower(first_name) LIKE ?")
		args = append(args, "%"+strings.ToLower(firstName)+"%")
	}
	if lastName != "" {
		query.WriteString(" AND lower(last_name) LIKE ?")
		args = append(args, "%"+strings.ToLower(lastName)+"%")
	}
    if firstName == "" && lastName == "" {
        // No criteria provided, return all or error based on desired behavior.
        // For now, returning all. Could also return an error or empty list.
        // return []models.Person{}, errors.New("search criteria (firstName or lastName) must be provided")
    }


	rows, err := r.DB.Query(query.String(), args...)
	if err != nil {
		log.Printf("Error querying for SearchByName: %v", err)
		return nil, err
	}
	defer rows.Close()

	persons := []models.Person{}
	for rows.Next() {
		var p models.Person
		err := rows.Scan(&p.SqliteID, &p.FirstName, &p.LastName, &p.Email, &p.Phone, &p.Address, &p.City, &p.State, &p.ZipCode)
		if err != nil {
			log.Printf("Error scanning row for SearchByName: %v", err)
			return nil, err
		}
		persons = append(persons, p)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error after iterating rows for SearchByName: %v", err)
		return nil, err
	}

	return persons, nil
}

// ListByCityState retrieves persons by city and state from SQLite.
func (r *SQLitePersonRepository) ListByCityState(city, state string) ([]models.Person, error) {
	if city == "" || state == "" {
		return nil, errors.New("city and state must be provided")
	}

	query := `SELECT id, first_name, last_name, email, phone, address, city, state, zip_code
	            FROM persons WHERE lower(city) = ? AND lower(state) = ?`

	rows, err := r.DB.Query(query, strings.ToLower(city), strings.ToLower(state))
	if err != nil {
		log.Printf("Error querying for ListByCityState: %v", err)
		return nil, err
	}
	defer rows.Close()

	persons := []models.Person{}
	for rows.Next() {
		var p models.Person
		err := rows.Scan(&p.SqliteID, &p.FirstName, &p.LastName, &p.Email, &p.Phone, &p.Address, &p.City, &p.State, &p.ZipCode)
		if err != nil {
			log.Printf("Error scanning row for ListByCityState: %v", err)
			// Consider returning partial results or a specific error
			return nil, err
		}
		persons = append(persons, p)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error after iterating rows for ListByCityState: %v", err)
		return nil, err
	}

	return persons, nil
}

// Helper to scan a single row into a Person model
func scanPerson(row *sql.Row) (*models.Person, error) {
    var p models.Person
    // Handle nullable fields by scanning into sql.NullString or similar, then converting.
    // For simplicity, assuming all text fields are non-null in DB or empty strings are acceptable.
    var phone, address, city, state, zipCode sql.NullString

    err := row.Scan(&p.SqliteID, &p.FirstName, &p.LastName, &p.Email, &phone, &address, &city, &state, &zipCode)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil // Not found
        }
        return nil, err // Other scan error
    }

    // Convert sql.NullString to string for the model
    p.Phone = phone.String
    p.Address = address.String
    p.City = city.String
    p.State = state.String
    p.ZipCode = zipCode.String

    return &p, nil
}

// Helper to scan multiple rows into a slice of Person models
func scanPersons(rows *sql.Rows) ([]models.Person, error) {
    persons := []models.Person{}
    for rows.Next() {
        var p models.Person
        var phone, address, city, state, zipCode sql.NullString // For nullable fields

        err := rows.Scan(&p.SqliteID, &p.FirstName, &p.LastName, &p.Email, &phone, &address, &city, &state, &zipCode)
        if err != nil {
            return nil, err // Error during row scan
        }

        // Convert sql.NullString to string
        p.Phone = phone.String
        p.Address = address.String
        p.City = city.String
        p.State = state.String
        p.ZipCode = zipCode.String

        persons = append(persons, p)
    }
    if err := rows.Err(); err != nil {
        return nil, err // Error after iterating rows (e.g., connection issue)
    }
    return persons, nil
}

// Update the GetPerson method to use scanPerson helper and handle nullable fields correctly
func (r *SQLitePersonRepository) GetPersonCorrected(idStr string) (*models.Person, error) {
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        return nil, fmt.Errorf("invalid ID format: %w", err)
    }

    query := `SELECT id, first_name, last_name, email, phone, address, city, state, zip_code FROM persons WHERE id = ?`
    row := r.DB.QueryRow(query, id)

    return scanPerson(row)
}

// Update SearchByName to use scanPersons helper
func (r *SQLitePersonRepository) SearchByNameCorrected(firstName, lastName string) ([]models.Person, error) {
    var queryBuilder strings.Builder
    queryBuilder.WriteString("SELECT id, first_name, last_name, email, phone, address, city, state, zip_code FROM persons WHERE 1=1")
    args := []interface{}{}

    if firstName != "" {
        queryBuilder.WriteString(" AND lower(first_name) LIKE ?")
        args = append(args, "%"+strings.ToLower(firstName)+"%")
    }
    if lastName != "" {
        queryBuilder.WriteString(" AND lower(last_name) LIKE ?")
        args = append(args, "%"+strings.ToLower(lastName)+"%")
    }

    rows, err := r.DB.Query(queryBuilder.String(), args...)
    if err != nil {
        return nil, fmt.Errorf("query error in SearchByName: %w", err)
    }
    defer rows.Close()

    return scanPersons(rows)
}


// Update ListByCityState to use scanPersons helper
func (r *SQLitePersonRepository) ListByCityStateCorrected(city, state string) ([]models.Person, error) {
    if city == "" || state == "" {
        return nil, errors.New("city and state must be provided")
    }

    query := `SELECT id, first_name, last_name, email, phone, address, city, state, zip_code
	            FROM persons WHERE lower(city) = ? AND lower(state) = ?`

    rows, err := r.DB.Query(query, strings.ToLower(city), strings.ToLower(state))
    if err != nil {
        return nil, fmt.Errorf("query error in ListByCityState: %w", err)
    }
    defer rows.Close()

    return scanPersons(rows)
}

// Replace original methods with corrected ones if these helpers are adopted.
// For now, the original methods are kept for directness, but these helpers show a way to handle nullable fields.
// The original GetPerson, SearchByName, ListByCityState already handle scanning directly.
// The main issue with direct scanning is that if a field is NULL in the DB, and you scan into a non-pointer string,
// it might error or convert NULL to an empty string, depending on the driver.
// sql.NullString correctly handles this by providing a Valid flag.

// Let's simplify and ensure the original methods correctly handle nullable fields by scanning into nullable types where appropriate.
// For GetPerson:
func (r *SQLitePersonRepository) GetPerson(idStr string) (*models.Person, error) { // Overwriting original
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("Invalid ID format: %s", idStr)
		return nil, errors.New("invalid ID format")
	}

	query := `SELECT id, first_name, last_name, email, phone, address, city, state, zip_code
	            FROM persons WHERE id = ?`
	row := r.DB.QueryRow(query, id)

	var p models.Person
    var phone, address, city, state, zipCode sql.NullString // Use sql.NullString for nullable fields

	err = row.Scan(&p.SqliteID, &p.FirstName, &p.LastName, &p.Email, &phone, &address, &city, &state, &zipCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Person not found
		}
		log.Printf("Error scanning person row: %v", err)
		return nil, err
	}
    // Assign to model, converting NullString to string
    p.Phone = phone.String
    p.Address = address.String
    p.City = city.String
    p.State = state.String
    p.ZipCode = zipCode.String

	return &p, nil
}

// For SearchByName:
func (r *SQLitePersonRepository) SearchByName(firstName, lastName string) ([]models.Person, error) { // Overwriting original
	var query strings.Builder
	query.WriteString("SELECT id, first_name, last_name, email, phone, address, city, state, zip_code FROM persons WHERE 1=1")
	args := []interface{}{}

	if firstName != "" {
		query.WriteString(" AND lower(first_name) LIKE ?")
		args = append(args, "%"+strings.ToLower(firstName)+"%")
	}
	if lastName != "" {
		query.WriteString(" AND lower(last_name) LIKE ?")
		args = append(args, "%"+strings.ToLower(lastName)+"%")
	}

	rows, err := r.DB.Query(query.String(), args...)
	if err != nil {
		log.Printf("Error querying for SearchByName: %v", err)
		return nil, err
	}
	defer rows.Close()

	persons := []models.Person{}
	for rows.Next() {
		var p models.Person
        var phone, address, city, state, zipCode sql.NullString
		err := rows.Scan(&p.SqliteID, &p.FirstName, &p.LastName, &p.Email, &phone, &address, &city, &state, &zipCode)
		if err != nil {
			log.Printf("Error scanning row for SearchByName: %v", err)
			return nil, err
		}
        p.Phone = phone.String
        p.Address = address.String
        p.City = city.String
        p.State = state.String
        p.ZipCode = zipCode.String
		persons = append(persons, p)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error after iterating rows for SearchByName: %v", err)
		return nil, err
	}

	return persons, nil
}

// For ListByCityState:
func (r *SQLitePersonRepository) ListByCityState(cityIn, stateIn string) ([]models.Person, error) { // Overwriting original
	if cityIn == "" || stateIn == "" {
		return nil, errors.New("city and state must be provided")
	}

	query := `SELECT id, first_name, last_name, email, phone, address, city, state, zip_code
	            FROM persons WHERE lower(city) = ? AND lower(state) = ?`

	rows, err := r.DB.Query(query, strings.ToLower(cityIn), strings.ToLower(stateIn))
	if err != nil {
		log.Printf("Error querying for ListByCityState: %v", err)
		return nil, err
	}
	defer rows.Close()

	persons := []models.Person{}
	for rows.Next() {
		var p models.Person
        var phone, address, city, state, zipCode sql.NullString
		err := rows.Scan(&p.SqliteID, &p.FirstName, &p.LastName, &p.Email, &phone, &address, &city, &state, &zipCode)
		if err != nil {
			log.Printf("Error scanning row for ListByCityState: %v", err)
			return nil, err
		}
        p.Phone = phone.String
        p.Address = address.String
        p.City = city.String
        p.State = state.String
        p.ZipCode = zipCode.String
		persons = append(persons, p)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error after iterating rows for ListByCityState: %v", err)
		return nil, err
	}

	return persons, nil
}
