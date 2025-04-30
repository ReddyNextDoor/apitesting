# Person API Service (Java Spring Boot)

A Spring Boot microservice for managing person records, supporting both SQLite and MongoDB backends. Includes CRUD, search, and flexible backend switching for local development, testing, Docker, and Kubernetes.

---

## 1. Project Overview
- **Framework:** Spring Boot + Springdoc OpenAPI
- **DB Backends:** SQLite (default), MongoDB
- **Features:** CRUD, search by name, list by city/state
- **Packaging:** Maven, Docker, Kubernetes

---

## 2. Directory Structure
- `src/main/java/com/example/personapi/` — Main application code
- `src/main/resources/` — Properties, static resources
- `db/` — SQLite DB files
- `Dockerfile`, `k8s.yaml` — Container & K8s deployment

---

## 3. Maven Setup

### Install Maven
See https://maven.apache.org/download.cgi

### Install Dependencies
```bash
mvn clean install
```

---

## 4. Running the Application

### SQLite (default)
```bash
PERSON_REPO_BACKEND=sqlite mvn spring-boot:run
```

### MongoDB
Ensure MongoDB is running locally or in Docker, then:
```bash
PERSON_REPO_BACKEND=mongo MONGO_URI="mongodb://user:pass@localhost:27017/person_db?authSource=admin" mvn spring-boot:run
```

- Access API docs at [http://localhost:8000/docs](http://localhost:8000/docs)

---

## 5. Running Tests

```bash
mvn test
```

---

## 6. Docker

### Build Image
```bash
docker build -t person-api-service-java .
```

### Run Container
```bash
docker run -d -p 8000:8000 \
  -e PERSON_REPO_BACKEND=sqlite \
  --name person-api-service-java person-api-service-java
```

- For MongoDB, add `-e PERSON_REPO_BACKEND=mongo -e MONGO_URI=...` as needed.

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

### Deploy (edit `k8s.yaml` as needed)
```bash
kubectl apply -f k8s.yaml
```

---

## 8. Environment Variables
- `PERSON_REPO_BACKEND`: `sqlite` (default) or `mongo`
- `MONGO_URI`: MongoDB connection string (if using mongo)
- `SQLITE_DB_PATH`: Path to SQLite DB file (default: `db/persons.db`)
