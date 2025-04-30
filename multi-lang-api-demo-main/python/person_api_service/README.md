# Person API Service

A FastAPI microservice for managing person records, supporting both SQLite and MongoDB backends. Includes CRUD, search, and flexible backend switching for local development, testing, Docker, and Kubernetes.

---

## 1. Project Overview
- **Framework:** FastAPI
- **DB Backends:** SQLite (default), MongoDB
- **Features:** CRUD, search by name, list by city/state
- **Packaging:** Poetry, Docker, Kubernetes

---

## 2. Directory Structure
- `app/` — Main FastAPI application code
- `app/models.py` — SQLAlchemy & Pydantic models
- `app/repository_sqlite.py` — SQLite repository
- `app/repository_mongo.py` — MongoDB repository
- `app/database.py` — DB engine/session setup
- `tests/` — Pytest test suite
- `db/` — SQLite DB files
- `Dockerfile`, `k8s.yaml` — Container & K8s deployment

---

## 3. Poetry Setup

### Install Poetry
```bash
curl -sSL https://install.python-poetry.org | python3 -
```

### Create/Activate Virtualenv
```bash
poetry env use python3.13
```

### Install Dependencies
```bash
poetry install
```

---

## 4. Running the Application

### SQLite (default)
```bash
PERSON_REPO_BACKEND=sqlite poetry run uvicorn app.main:app --reload --port 8000
```

### MongoDB
Ensure MongoDB is running locally or in Docker, then:
```bash
PERSON_REPO_BACKEND=mongo MONGO_URI="mongodb://user:pass@localhost:27017/person_db?authSource=admin" poetry run uvicorn app.main:app --reload --port 8000
```

- Access API docs at [http://localhost:8000/docs](http://localhost:8000/docs)

---

## 5. Running Tests

### SQLite (default)
```bash
poetry run pytest --backend=sqlite -s
```

### MongoDB
```bash
poetry run pytest --backend=mongo -s
```

---

## 6. Docker

### Build Image
```bash
docker build -t person-api-service .
```

### Run Container
```bash
docker run -d -p 8000:8000 \
  -e PERSON_REPO_BACKEND=sqlite \
  --name person-api-service person-api-service
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

- The API will be available at [http://localhost:8000/docs](http://localhost:8000/docs)

---

## 7. Kubernetes

### Deploy All Services
Apply the manifest to create both SQLite and MongoDB API pods (and MongoDB itself):
```bash
kubectl apply -f k8s.yaml
```

### Access the API (choose backend)
- **SQLite backend:**
  - The API is available at [http://localhost:30080/docs](http://localhost:30080/docs)
- **MongoDB backend:**
  - The API is available at [http://localhost:30081/docs](http://localhost:30081/docs)

If running on a remote cluster, use the appropriate node's external IP and the corresponding port.

- Both API pods are always available; access the desired backend via its NodePort.
- MongoDB itself is only exposed as a ClusterIP for internal use by the API pod.

---

## 8. Environment Variables

| Variable             | Description                              | Example                                      |
|----------------------|------------------------------------------|----------------------------------------------|
| PERSON_REPO_BACKEND  | Backend to use: `sqlite` or `mongo`      | `sqlite`                                     |
| SQLITE_DB_PATH       | SQLite DB file path                      | `db/persons.db` or `db/test_persons.db`      |
| MONGO_URI            | MongoDB connection string                | `mongodb://user:pass@localhost:27017/...`    |
| MONGO_DB             | MongoDB DB name (for tests)              | `test_person_db`                             |

---

- Ensure `db/` exists for SQLite file-based DBs.
- For MongoDB tests, ensure a local or Docker MongoDB instance is running.
- All tests/services support backend switching via environment variables or pytest CLI.
