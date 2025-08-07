# apitesting

A multi-language REST API demonstration repository.

## Overview

This repository showcases REST API implementations in four different programming languages:

- **Python** (FastAPI)
- **TypeScript** (Node.js/Express)
- **Java** (Spring Boot)
- **C#** (.NET Core)

Each subfolder under `multi-lang-api-demo-main/` contains an independent, fully functional REST API service for managing "Person" resources, demonstrating similar features and endpoints in each language.

---

## Directory Structure

```
multi-lang-api-demo-main/
  python/       # FastAPI implementation
  typescript/   # Node.js (Express) implementation
  java/         # Spring Boot implementation
  dotnet/       # ASP.NET Core implementation
LICENSE
```

---

## Functionality

Each API provides the following core features:
- CRUD operations on a `Person` resource (create, read, update, delete)
- Search/filter capabilities (e.g., search by name, list by city/state)
- Swagger/OpenAPI documentation (where supported)
- Dockerfile for containerized deployment

### Example Endpoints

- `GET /api/person` — List all persons
- `GET /api/person/{id}` — Get a person by ID
- `POST /api/person` — Create a new person
- `PUT /api/person/{id}` — Update a person
- `DELETE /api/person/{id}` — Delete a person
- `GET /api/person/search?firstName=X&lastName=Y` — Search by name
- `GET /api/person/by_city_state?city=X&state=Y` — List by city/state

---

## Getting Started

### Prerequisites

- Docker (recommended for all stacks)
- [Python 3.8+](https://www.python.org/) (for FastAPI)
- [Node.js 16+](https://nodejs.org/) (for TypeScript/Express)
- [Java 17+](https://adoptopenjdk.net/) (for Spring Boot)
- [.NET 9 SDK](https://dotnet.microsoft.com/) (for ASP.NET Core)

### Running with Docker

Navigate to any implementation and build/run the Docker image:

```sh
cd multi-lang-api-demo-main/<language>
docker build -t <api-name> .
docker run -p 8000:8000 <api-name>
```

Replace `<language>` and `<api-name>` with your chosen implementation (e.g., `python`, `typescript`, `java`, `dotnet`).

### Local Development

See each subfolder's code and documentation for local (non-Docker) instructions. Most support running via language-native tools (e.g., `uvicorn`, `npm start`, `mvn spring-boot:run`, `dotnet run`).

---

## API Documentation

- Swagger/OpenAPI documentation is available at `/docs` or `/swagger` on each service (see individual subfolders).

---

## License

MIT License — see [LICENSE](LICENSE) for details.

---

## Contributing

Pull requests and issues are welcome. Please ensure your code adheres to the conventions and standards of the respective language implementation.
