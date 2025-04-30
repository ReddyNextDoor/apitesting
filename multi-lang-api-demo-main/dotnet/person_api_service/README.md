# Person API Service (.NET Core)

A .NET Core Web API for managing person records, supporting both SQLite and MongoDB backends. Includes CRUD, search, and flexible backend switching for local development, testing, Docker, and Kubernetes.

---

## 1. Project Overview
- **Framework:** ASP.NET Core
- **DB Backends:** SQLite (default), MongoDB
- **Features:** CRUD, search by name, list by city/state
- **Packaging:** dotnet CLI, Docker, Kubernetes

---

## 2. Directory Structure
- `Controllers/` — API controllers
- `Models/` — Entity and DTO classes
- `Repositories/` — Repository interfaces and implementations
- `db/` — SQLite DB files
- `Dockerfile`, `k8s.yaml` — Container & K8s deployment

---

## 3. .NET CLI Setup

### Install .NET SDK
See https://dotnet.microsoft.com/download

### Restore Dependencies
```bash
dotnet restore
```

---

## 4. Running the Application

### SQLite (default)
```bash
PERSON_REPO_BACKEND=sqlite dotnet run
```

### MongoDB

Ensure MongoDB is running in your environment (locally, in Docker, or in Kubernetes).

#### Local Example (custom port):
If your MongoDB is running locally on port 27080, use:
```bash
dotnet build
MONGO_URI="mongodb://superuser:supersecret@localhost:27080/person_db?authSource=admin" PERSON_REPO_BACKEND=mongo dotnet run
```

#### Docker Example:
If using Docker Compose and your MongoDB service is named `mongo` (default port 27017):
```bash
PERSON_REPO_BACKEND=mongo MONGO_URI="mongodb://superuser:supersecret@mongo:27017/person_db?authSource=admin" dotnet run
```

#### Kubernetes Example:
Set the environment variables in your deployment YAML:
```yaml
- name: PERSON_REPO_BACKEND
  value: "mongo"
- name: MONGO_URI
  value: "mongodb://superuser:supersecret@mongo:27017/person_db?authSource=admin"
```

> **Note:** The `MONGO_URI` environment variable must always be set explicitly when using the mongo backend. The application will fail to start if it is not set.

- Access API docs at [http://localhost:8000/docs](http://localhost:8000/docs)

---

## 5. Running Tests

```bash
dotnet test
```

---

## 6. Docker

### Build the Docker Image
```bash
docker build -t person-api-service-dotnet .
```

### Run the Container (SQLite backend)
```bash
docker run -d \
  -p 8000:8000 \
  -e PERSON_REPO_BACKEND=sqlite \
  --name person-api-service-dotnet person-api-service-dotnet
```

### Run the Container (MongoDB backend)
Ensure MongoDB is running and accessible to the container. Pass the correct MONGO_URI:
```bash
docker run -d \
  -p 8000:8000 \
  -e PERSON_REPO_BACKEND=mongo \
  -e MONGO_URI="mongodb://superuser:supersecret@mongo:27017/person_db?authSource=admin" \
  --name person-api-service-dotnet-mongo person-api-service-dotnet
```

> **Tip:** For local development with Docker Compose, ensure your MongoDB service is named `mongo` and exposes port 27017.

### Docker Compose (Recommended for Multi-Service)

#### SQLite Backend
```bash
docker-compose --profile sqlite up --build
```

#### MongoDB Backend
```bash
docker-compose --profile mongo up --build
```

---

## 7. Kubernetes

### Build the Docker Image (if running locally/minikube)
```bash
docker build -t person-api-service-dotnet:latest .
```

### Deploy to Kubernetes

Edit `k8s.yaml` as needed for your environment. The manifest already includes deployments and services for both SQLite and MongoDB backends.

Apply the manifest:
```bash
kubectl apply -f k8s.yaml
```

- For MongoDB backend, ensure a MongoDB service is accessible at the URI specified in the deployment's `MONGO_URI`.
- You may need to load the local Docker image into your Kubernetes cluster (e.g., with minikube: `minikube image load person-api-service-dotnet:latest`).

Access the API via the NodePort specified in the manifest (default: 30080 for SQLite, 30081 for MongoDB).

---

## 8. Environment Variables
- `PERSON_REPO_BACKEND`: `sqlite` (default) or `mongo`
- `MONGO_URI`: MongoDB connection string (if using mongo)
- `SQLITE_DB_PATH`: Path to SQLite DB file (default: `db/persons.db`)
