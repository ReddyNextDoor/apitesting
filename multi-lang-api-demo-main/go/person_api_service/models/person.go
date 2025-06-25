package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv" // Moved import here
)

// Person represents the main data model for a person.
// For MongoDB, `bson:"_id,omitempty"` is used for the ID field.
// For SQLite, `gorm:"primaryKey"` would be used if GORM were the ORM.
// Since we're using raw SQL for SQLite for simplicity, ID handling will be manual there.
type Person struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"` // MongoDB uses ObjectID
	SqliteID  int64              `json:"sqlite_id,omitempty" gorm:"primaryKey" bson:"-"` // SQLite uses auto-incrementing INT
	FirstName string             `json:"first_name" binding:"required"`
	LastName  string             `json:"last_name" binding:"required"`
	Email     string             `json:"email" binding:"required,email"`
	Phone     string             `json:"phone,omitempty"`
	Address   string             `json:"address,omitempty"`
	City      string             `json:"city,omitempty"`
	State     string             `json:"state,omitempty"`
	ZipCode   string             `json:"zip_code,omitempty"`
}

// PersonCreate is used when creating a new person. ID is omitted.
type PersonCreate struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone,omitempty"`
	Address   string `json:"address,omitempty"`
	City      string `json:"city,omitempty"`
	State     string `json:"state,omitempty"`
	ZipCode   string `json:"zip_code,omitempty"`
}

// PersonUpdate is used when updating an existing person. All fields are optional.
type PersonUpdate struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Address   string `json:"address,omitempty"`
	City      string `json:"city,omitempty"`
	State     string `json:"state,omitempty"`
	ZipCode   string `json:"zip_code,omitempty"`
}

// PersonOut is used for API responses. It includes the ID.
// This struct can be the same as Person for now, or customized if needed.
type PersonOut struct {
	ID        string `json:"id"` // Represent ID as string in output for both DBs
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone,omitempty"`
	Address   string `json:"address,omitempty"`
	City      string `json:"city,omitempty"`
	State     string `json:"state,omitempty"`
	ZipCode   string `json:"zip_code,omitempty"`
}

// ToPersonOut converts a Person (from DB) to PersonOut (for API response).
// It handles the ID conversion based on which ID is populated.
func (p *Person) ToPersonOut() PersonOut {
	var idStr string
	if !p.ID.IsZero() { // Check if MongoDB ObjectID is populated
		idStr = p.ID.Hex()
	} else if p.SqliteID != 0 { // Check if SQLite ID is populated
		idStr = strconv.FormatInt(p.SqliteID, 10)
	}

	return PersonOut{
		ID:        idStr,
		FirstName: p.FirstName,
		LastName:  p.LastName,
		Email:     p.Email,
		Phone:     p.Phone,
		Address:   p.Address,
		City:      p.City,
		State:     p.State,
		ZipCode:   p.ZipCode,
	}
}

// Helper function for converting PersonCreate to Person
func (pc *PersonCreate) ToPerson() Person {
	return Person{
		FirstName: pc.FirstName,
		LastName:  pc.LastName,
		Email:     pc.Email,
		Phone:     pc.Phone,
		Address:   pc.Address,
		City:      pc.City,
		State:     pc.State,
		ZipCode:   pc.ZipCode,
	}
}

// This is needed for the ToPersonOut function
// import "strconv" // Already imported at the top
