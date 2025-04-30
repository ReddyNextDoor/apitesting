# Person API Service (TypeScript/Node.js)

A Node.js REST API for managing person records, supporting both SQLite and MongoDB backends. Includes CRUD, search, and flexible backend switching for local development, Docker Compose, and Kubernetes.

---

## 1. Project Overview
- **Framework:** Express.js + TypeScript
- **DB Backends:** SQLite (default, via TypeORM), MongoDB (via Mongoose)
- **Features:** CRUD, search by name, list by city/state
- **Packaging:** npm scripts, Docker, Kubernetes

---

## 2. Directory Structure
- `src/` — Source code (models, repositories, controllers, config)
- `db/` — SQLite DB files
- `tests/` — Jest-based tests
- `Dockerfile`, `docker-compose.yml`, `k8s.yaml`, `openapi.yaml` — Container & K8s deployment

---

## 3. Local Development

### Install dependencies
```bash
npm install
```

### Build
```bash
npm run build
```

### Run (SQLite, default)
```bash
npm run dev
```

### Run (MongoDB)
Ensure MongoDB is running locally or in Docker, then:
```bash
PERSON_REPO_BACKEND=mongo MONGO_URI="mongodb://superuser:supersecret@localhost:27017/person_db?authSource=admin" npm run dev
```

- Access API docs at [http://localhost:8000/docs](http://localhost:8000/docs)

---

## 4. Docker Compose

### Build and Run (SQLite backend)
```bash
docker compose --profile sqlite up --build
```
- SQLite API available at [http://localhost:8000/docs](http://localhost:8000/docs)

### Build and Run (MongoDB backend)
```bash
docker compose --profile mongo up --build
```
- MongoDB and API available at [http://localhost:8000/docs](http://localhost:8000/docs)

---

## 5. Kubernetes

### Build Docker Image (for local clusters)
```bash
docker build -t person-api-service-typescript:latest .
```

### (Optional) Push to Registry (for remote clusters)
```bash
docker tag person-api-service-typescript:latest <your-dockerhub-username>/person-api-service-typescript:latest
docker push <your-dockerhub-username>/person-api-service-typescript:latest
```
Edit `k8s.yaml` to use your registry image if needed.

### Deploy to Kubernetes
```bash
kubectl apply -f k8s.yaml
```

#### Access the API
- **SQLite:** [http://localhost:30080/docs](http://localhost:30080/docs)
- **MongoDB:** [http://localhost:30081/docs](http://localhost:30081/docs)

To access the API from your local machine, use kubectl port-forward:
```bash
kubectl port-forward service/ts-api-sqlite 8000:8000
kubectl port-forward service/ts-api-mongo 8000:8000
```

Alternatively, if your cluster exposes NodePort services directly, access via:
- SQLite: http://localhost:30080/docs
- MongoDB: http://localhost:30081/docs

---

## 6. Troubleshooting

- If Swagger UI points to the wrong port, ensure the `servers:` field in `openapi.yaml` is set to `url: /`.
- If NodePort endpoints return "Page Not Found", check pod logs:
  ```bash
  kubectl get pods
  kubectl logs <pod-name>
  ```
- For persistent MongoDB data, update `k8s.yaml` to use a PersistentVolumeClaim instead of `emptyDir`.

---

## 7. Testing

### Run tests
```bash
npm test
```

---

## 8. API Documentation

- OpenAPI spec: [`openapi.yaml`](./openapi.yaml)
- Interactive docs: `/docs` endpoint when running

---

## 9. Environment Variables

- `PERSON_REPO_BACKEND`: `sqlite` (default) or `mongo`
- `MONGO_URI`: MongoDB connection string (required for MongoDB backend)

---

## 10. Notes
- Ensure only one API service binds to a port at a time (avoid port conflicts).
- For production, use a real database volume instead of `emptyDir` in Kubernetes.

```

### Run the container (MongoDB backend)
```bash
docker run -d \
  -p 8000:8000 \
  -e PERSON_REPO_BACKEND=mongo \
  -e MONGO_URI="mongodb://user:pass@mongo:27017/person_db?authSource=admin" \
  --name person-api-service-typescript-mongo person-api-service-typescript
```

---

## 5. Kubernetes

### Build the Docker image (if running locally/minikube)
```bash
docker build -t person-api-service-typescript:latest .
```

### Deploy to Kubernetes
```bash
kubectl apply -f k8s.yaml
```

- Edit `k8s.yaml` as needed for your environment.
- Access the API via the NodePort specified in the manifest.

---

## 6. Environment Variables
- `PERSON_REPO_BACKEND`: `sqlite` (default) or `mongo`
- `SQLITE_DB_PATH`: Path to SQLite DB file (default: `db/persons.db`)
- `MONGO_URI`: MongoDB connection string (required if backend is mongo)

---

## 7. Testing

### Run tests
```bash
npm test
```

---

## 8. API Documentation
- Swagger UI available at `/docs`
- OpenAPI schema available at `/openapi.json`

---

## 9. License
MIT
