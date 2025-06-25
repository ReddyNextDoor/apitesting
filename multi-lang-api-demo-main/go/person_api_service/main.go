package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/example/person_api_service/database"
	"github.com/example/person_api_service/handlers"
	"github.com/gin-gonic/gin"

	// Swagger imports for documentation (if using swaggo)
	// _ "github.com/example/person_api_service/docs" // If you generate docs
	// ginSwagger "github.com/swaggo/gin-swagger"
	// swaggerFiles "github.com/swaggo/files"
)

// @title Person API Service (Go)
// @version 1.0
// @description This is a Go implementation of the Person API service.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http https
func main() {
	// Determine backend and initialize databases
	backend := strings.ToLower(os.Getenv("PERSON_REPO_BACKEND"))
	log.Printf("PERSON_REPO_BACKEND: %s", backend)

	if backend == "mongo" {
		log.Println("Initializing MongoDB...")
		database.InitMongoDB()
		defer database.CloseMongoDB()
	} else { // Default to SQLite
		log.Println("Initializing SQLite...")
		database.InitSQLite()
		defer database.CloseSqliteDB()
	}

	// Set Gin mode (e.g., release, debug, test)
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.ReleaseMode
	}
	gin.SetMode(ginMode)
	log.Printf("GIN_MODE: %s", ginMode)


	router := gin.Default()

	// Middleware for logging (Gin's default logger is quite good)
	// router.Use(gin.Logger())
	// Middleware for recovery
	router.Use(gin.Recovery())

	// Simple redirect from root to /docs (if using Swagger) or a welcome message
	router.GET("/", func(c *gin.Context) {
		// For now, let's redirect to /health as a basic root handler
		// Later, if Swagger UI is added, this can redirect to /swagger/index.html
		c.Redirect(http.StatusMovedPermanently, "/health")
	})


	// Initialize handlers
	personHandler := handlers.NewPersonHandler()

	// Health check endpoint
	router.GET("/health", handlers.HealthCheck)

	// Person API routes
	personRoutes := router.Group("/persons")
	{
		personRoutes.POST("/", personHandler.CreatePerson)
		personRoutes.GET("/search", personHandler.SearchPersons)
		personRoutes.GET("/by_city_state", personHandler.PersonsByCityState)
		personRoutes.GET("/:person_id", personHandler.GetPerson)
		personRoutes.PUT("/:person_id", personHandler.UpdatePerson)
		personRoutes.DELETE("/:person_id", personHandler.DeletePerson)
	}

	// Swagger documentation endpoint (if using swaggo)
	// Make sure to generate docs first: `swag init` in the service directory
	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// log.Println("Swagger UI available at /swagger/index.html")


	// Server configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Server starting on port %s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", port, err)
	}

	log.Println("Server stopped")
}
