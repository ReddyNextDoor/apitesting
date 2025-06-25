# Person API Service (Go Implementation)

A Gin-based microservice for managing person records, supporting both SQLite and MongoDB backends. This Go version provides CRUD operations, search functionality, and flexible backend switching for local development, testing, Docker, and Kubernetes deployments.

---

## 1. Project Overview
- **Framework:** Gin Web Framework
- **DB Backends:** SQLite (default), MongoDB
- **Features:**
    - CRUD operations for person records (Create, Read, Update, Delete)
    - Search persons by first name and/or last name
    - List persons by city and state
    - Health check endpoint
- **Packaging:** Go Modules, Docker, Docker Compose, Kubernetes

---

## 2. Directory Structure
```
.
├── Dockerfile             # For building the Docker image
├── docker-compose.yml     # For running with Docker Compose (SQLite & MongoDB profiles)
├── go.mod                 # Go module definition
├── go.sum                 # Go module checksums
├── k8s.yaml               # Kubernetes deployment manifest
├── README.md              # This file
├── db/                    # Default directory for SQLite database files (created automatically)
├── database/              # Database connection logic (SQLite, MongoDB)
│   └── database.go
├── handlers/              # HTTP request handlers (controllers)
│   └── person_handlers.go
├── main.go                # Main application entry point
├── models/                # Data structures (Person model, DTOs)
│   └── person.go
└── repository/            # Data access layer (interfaces and implementations)
    ├── mongo_repository.go
    ├── repository_interface.go
    └── sqlite_repository.go
    └── (tests for repository will also be here, e.g., sqlite_repository_test.go)
```

---

## 3. Prerequisites
- **Go:** Version 1.21 or higher.
- **Docker:** Latest version.
- **Docker Compose:** Latest version.
- **kubectl:** For Kubernetes deployment (connected to a K8s cluster like Minikube, Kind, or a cloud provider).
- **MongoDB Server** (Optional, for local MongoDB backend testing without Docker Compose): Install from [MongoDB website](https://www.mongodb.com/try/download/community).

---

## 4. Setup

1.  **Clone the repository (if you haven't already):**
    ```bash
    # Assuming you are in the root of the multi-lang-api-demo
    # This service is typically part of a larger project structure.
    ```

2.  **Navigate to the Go service directory:**
    ```bash
    cd multi-lang-api-demo-main/go/person_api_service
    ```

3.  **Install dependencies:**
    Go modules should handle this automatically when you build or run. You can explicitly download them:
    ```bash
    go mod download
    go mod tidy
    ```

---

## 5. Running the Application Locally

The application listens on port `8080` by default.

### SQLite Backend (Default)
This is the default if `PERSON_REPO_BACKEND` is not set or set to `sqlite`.
```bash
go run main.go
# or to specify backend explicitly:
# PERSON_REPO_BACKEND=sqlite go run main.go
```
The SQLite database file will be created at `db/persons.db` by default.

### MongoDB Backend
1.  **Ensure MongoDB is running:**
    You can start a local MongoDB instance or use Docker:
    ```bash
    # Using Docker (if you don't have a local MongoDB server)
    docker run -d -p 27017:27017 --name local-mongo mongo:6
    ```
2.  **Run the application with MongoDB environment variables:**
    ```bash
    PERSON_REPO_BACKEND=mongo \
    MONGO_URI="mongodb://localhost:27017/local_person_db_go" \
    MONGO_DB="local_person_db_go" \
    go run main.go
    ```

-   **API Access:**
    -   Health Check: [http://localhost:8080/health](http://localhost:8080/health)
    -   Other endpoints are available under `/persons/`. (e.g., `GET /persons/search?first_name=John`)
    -   Swagger/OpenAPI docs might be available at `/swagger/index.html` if integrated and generated.

---

## 6. Running Tests
Execute unit tests for the repository and other components:
```bash
# From the go/person_api_service directory
go test ./... -v
```
This will run all `*_test.go` files in the current directory and its subdirectories.
The SQLite repository tests will create and use a temporary database file in `test_db_files/`.

---

## 7. Docker

### Build the Docker Image
```bash
# From the go/person_api_service directory
docker build -t go-person-api-service:latest .
```

### Run the Container

**With SQLite Backend:**
```bash
docker run -d -p 8080:8080 \
  -e PERSON_REPO_BACKEND=sqlite \
  -e SQLITE_DB_PATH=db/docker_persons.db \
  --name go-api-sqlite go-person-api-service:latest
```
Access at: `http://localhost:8080/health`

**With MongoDB Backend (requires a running MongoDB instance accessible to Docker):**
If you have a MongoDB container named `my_mongo_container` on the same Docker network:
```bash
docker run -d -p 8081:8080 \
  -e PERSON_REPO_BACKEND=mongo \
  -e MONGO_URI="mongodb://my_mongo_container:27017/docker_person_db" \
  -e MONGO_DB="docker_person_db" \
  --name go-api-mongo go-person-api-service:latest
  # --network <your_docker_network> # If mongo is on a specific network
```
Access at: `http://localhost:8081/health` (Note the different host port `8081`)

### Docker Compose (Recommended for local multi-service setups)

**SQLite Backend:**
```bash
# From the go/person_api_service directory
docker-compose --profile sqlite up --build
```
API available at: [http://localhost:8080/health](http://localhost:8080/health)
SQLite data is persisted in a Docker volume named `db_compose_sqlite`.

**MongoDB Backend:**
This profile also starts a MongoDB container.
```bash
# From the go/person_api_service directory
docker-compose --profile mongo up --build
```
API available at: [http://localhost:8081/health](http://localhost:8081/health)
MongoDB data is persisted in a Docker volume named `mongo_data_go`.
MongoDB is accessible on host port `27018` (e.g., `mongodb://localhost:27018/`).

---

## 8. Kubernetes

The `k8s.yaml` manifest deploys:
- A MongoDB instance (with PVC for persistence).
- The Go Person API service connected to SQLite (NodePort `30090`).
- The Go Person API service connected to the deployed MongoDB (NodePort `30091`).

1.  **Ensure your Docker image is accessible to Kubernetes:**
    - If using Minikube/Kind, you might load the local image:
      ```bash
      # Example for Minikube
      # eval $(minikube -p <profile_name> docker-env) # Point Docker CLI to Minikube's daemon
      # docker build -t go-person-api-service:latest .
      # Or: minikube image load go-person-api-service:latest --profile <profile_name>

      # Example for Kind
      # kind load docker-image go-person-api-service:latest --name <kind_cluster_name>
      ```
    - Or, push the image to a container registry (Docker Hub, GCR, ECR, etc.) and update `image:` in `k8s.yaml` if needed. The default `imagePullPolicy` is `IfNotPresent` if tag is `latest`, otherwise `Always`.

2.  **Deploy to Kubernetes:**
    ```bash
    # From the go/person_api_service directory
    kubectl apply -f k8s.yaml
    ```

3.  **Access the API:**
    - **SQLite backend:** `http://<K8S_NODE_IP>:30090/health`
    - **MongoDB backend:** `http://<K8S_NODE_IP>:30091/health`
    Replace `<K8S_NODE_IP>` with the IP address of one of your Kubernetes nodes. For Minikube, you can get it via `minikube ip`.

    To find services and their NodePorts:
    ```bash
    kubectl get services
    ```

---

## 9. Environment Variables

| Variable            | Description                                       | Default (in Dockerfile) / Example             |
|---------------------|---------------------------------------------------|-----------------------------------------------|
| `GIN_MODE`          | Gin framework mode                                | `release`                                     |
| `PORT`              | Port the application listens on inside container  | `8080`                                        |
| `PERSON_REPO_BACKEND`| Backend to use: `sqlite` or `mongo`             | `sqlite`                                      |
| `SQLITE_DB_PATH`    | SQLite database file path (relative to workdir)   | `db/persons.db`                               |
| `MONGO_URI`         | MongoDB connection string                         | `mongodb://mongo:27017/person_db` (for Docker) / `mongodb://mongo-service-go:27017/person_k8s_db_go` (for K8s) |
| `MONGO_DB`          | MongoDB database name                             | `person_db` (for Docker) / `person_k8s_db_go` (for K8s) |

---

## 10. API Endpoints

Base URL: `http://localhost:8080` (or relevant host/port for Docker/K8s)

-   **Health Check:**
    -   `GET /health`
-   **Persons:**
    -   `POST /persons/` - Create a new person.
        -   Request Body: `models.PersonCreate` JSON
    -   `GET /persons/search?first_name=<fn>&last_name=<ln>` - Search persons.
    -   `GET /persons/by_city_state?city=<city>&state=<state>` - List persons by city/state.
    -   `GET /persons/{person_id}` - Get a person by ID.
    -   `PUT /persons/{person_id}` - Update a person by ID.
        -   Request Body: `models.PersonUpdate` JSON
    -   `DELETE /persons/{person_id}` - Delete a person by ID.

Refer to `models/person.go` for `PersonCreate`, `PersonUpdate`, and `PersonOut` structures.

---
