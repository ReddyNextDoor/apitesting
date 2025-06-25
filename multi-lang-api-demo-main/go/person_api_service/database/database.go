package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	// MongoClient is the global MongoDB client.
	MongoClient *mongo.Client
	// SqliteDB is the global SQLite database connection.
	SqliteDB *sql.DB
)

// InitSQLite initializes the SQLite database connection and creates the table if it doesn't exist.
func InitSQLite() {
	dbPath := os.Getenv("SQLITE_DB_PATH")
	if dbPath == "" {
		dbPath = "db/persons.db" // Default path
	}

	// Ensure the db directory exists
	dbDir := "db"
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			log.Fatalf("Failed to create database directory: %v", err)
		}
	}


	var err error
	SqliteDB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open SQLite database: %v", err)
	}

	// Create persons table if it doesn't exist
	createTableSQL := `CREATE TABLE IF NOT EXISTS persons (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"first_name" TEXT NOT NULL,
		"last_name" TEXT NOT NULL,
		"email" TEXT NOT NULL UNIQUE,
		"phone" TEXT,
		"address" TEXT,
		"city" TEXT,
		"state" TEXT,
		"zip_code" TEXT
	);`
	_, err = SqliteDB.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create persons table: %v", err)
	}
	log.Println("SQLite database initialized and table created successfully.")
}

// InitMongoDB initializes the MongoDB client.
func InitMongoDB() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017" // Default URI
		log.Println("MONGO_URI not set, using default:", mongoURI)
	}

	var err error
	clientOptions := options.Client().ApplyURI(mongoURI)
	MongoClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping the primary
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := MongoClient.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	log.Println("Connected to MongoDB successfully!")
}

// GetMongoDB returns the "person_db" database from the MongoDB client.
// The actual database name can be made configurable via an environment variable if needed.
func GetMongoDB() *mongo.Database {
	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		dbName = "person_db" // Default database name
	}
	return MongoClient.Database(dbName)
}

// CloseMongoDB closes the MongoDB client connection.
func CloseMongoDB() {
	if MongoClient != nil {
		if err := MongoClient.Disconnect(context.TODO()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		} else {
			log.Println("MongoDB connection closed.")
		}
	}
}

// CloseSqliteDB closes the SQLite database connection.
func CloseSqliteDB() {
	if SqliteDB != nil {
		if err := SqliteDB.Close(); err != nil {
			log.Printf("Error closing SQLite DB: %v", err)
		} else {
			log.Println("SQLite DB connection closed.")
		}
	}
}

// GetPersonCollection is a helper to get the persons collection from MongoDB.
func GetPersonCollection() *mongo.Collection {
	return GetMongoDB().Collection("persons")
}

// GetDB returns the appropriate database client based on the environment variable.
// This is a conceptual function; actual repository implementations will use SqliteDB or MongoClient directly.
func GetDB() interface{} {
	backend := os.Getenv("PERSON_REPO_BACKEND")
	if backend == "mongo" {
		if MongoClient == nil {
			InitMongoDB()
		}
		return GetMongoDB()
	}
	if SqliteDB == nil {
		InitSQLite()
	}
	return SqliteDB
}
