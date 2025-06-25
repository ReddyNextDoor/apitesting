# Stage 1: Build the Go application
FROM golang:1.21-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download && go mod verify

# Copy the source code into the container
COPY . .

# Build the Go app
# -ldflags="-w -s" reduces the size of the binary by removing debug information
# CGO_ENABLED=0 ensures static linking for alpine compatibility if not using netgo
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-w -s" -o /go/bin/person-api-service ./main.go

# Stage 2: Create the final lightweight image
FROM alpine:latest

# Add ca-certificates for HTTPS calls if needed by the application
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/bin/person-api-service .

# Create a non-root user and group for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Expose port 8080 to the outside world
EXPOSE 8080

# Environment variables (can be overridden at runtime)
ENV GIN_MODE=release
ENV PORT=8080
# Default to SQLite backend. Can be set to "mongo" to use MongoDB.
ENV PERSON_REPO_BACKEND=sqlite
# Default SQLite database path (relative to workdir inside container, if app creates it)
# For SQLite, the db file will be created inside the container if it doesn't exist.
# Consider mounting a volume for persistent SQLite data in production.
ENV SQLITE_DB_PATH=db/persons.db
# Default MongoDB URI (points to a mongo service, typically on the Docker network or K8s cluster)
ENV MONGO_URI="mongodb://mongo:27017/person_db"
# Default MongoDB Database name
ENV MONGO_DB="person_db"


# Command to run the executable
# The db directory will be created by the application if SQLITE_DB_PATH includes a subdirectory
# and the SQLite initialization logic handles directory creation.
# Our InitSQLite function in database/database.go does create the "db" directory.
CMD ["./person-api-service"]
