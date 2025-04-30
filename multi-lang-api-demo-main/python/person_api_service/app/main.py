import os
from typing import List, Optional

from fastapi import Depends, FastAPI, HTTPException

from app import database, models, repository_sqlite
from app.models import PersonCreate, PersonOut, PersonUpdate
from app.repository_interface import PersonRepositoryInterface
from app.repository_mongo import MongoPersonRepository

models.Base.metadata.create_all(bind=database.engine)

from fastapi.responses import RedirectResponse

from app.openapi_config import custom_openapi

app = FastAPI()

app.openapi = lambda: custom_openapi(app)

@app.get("/", include_in_schema=False)
def root():
    return RedirectResponse(url="/docs")

@app.get("/docs/index.html", include_in_schema=False)
def docs_index_redirect():
    return RedirectResponse(url="/docs")

@app.get("/health")
def health():
    backend = os.getenv("PERSON_REPO_BACKEND", "sqlite").lower()
    mongo_db = os.getenv("MONGO_DB", "person_db")
    sqlite_db_path = os.getenv("SQLITE_DB_PATH", "db/persons.db")
    return {
        "backend": backend,
        "mongo_db": mongo_db,
        "sqlite_db_path": sqlite_db_path
    }

# Dependency for SQLite

def get_sqlite_repo():
    db = database.SessionLocal()
    try:
        yield repository_sqlite.PersonRepository(db)
    finally:
        db.close()

# Dependency for MongoDB

def get_mongo_repo():
    yield MongoPersonRepository()

# Select repository implementation based on environment variable
from loguru import logger

def get_repo():
    backend = os.getenv("PERSON_REPO_BACKEND", "sqlite").lower()
    if backend == "mongo":
        logger.info(f"[main.py] Using Mongo backend. DB: {os.getenv('MONGO_DB', 'person_db')}")
        return next(get_mongo_repo())
    else:
        logger.info(f"[main.py] Using SQLite backend. DB file: {os.getenv('SQLITE_DB_PATH', 'db/persons.db')}")
        return next(get_sqlite_repo())

# Log at startup
logger.info(f"[main.py] PERSON_REPO_BACKEND={os.getenv('PERSON_REPO_BACKEND', 'sqlite')}")
logger.info(f"[main.py] MONGO_DB={os.getenv('MONGO_DB', 'person_db')}")
logger.info(f"[main.py] SQLITE_DB_PATH={os.getenv('SQLITE_DB_PATH', 'db/persons.db')}")

@app.post("/persons/", response_model=PersonOut)
def create_person(person: PersonCreate, repo: PersonRepositoryInterface = Depends(get_repo)):
    return repo.create_person(person)

@app.get("/persons/search", response_model=List[PersonOut])
def search_persons(
    first_name: Optional[str] = None,
    last_name: Optional[str] = None,
    repo: PersonRepositoryInterface = Depends(get_repo),
):
    return repo.search_by_name(first_name, last_name)

@app.get("/persons/by_city_state", response_model=List[PersonOut])
def persons_by_city_state(
    city: str,
    state: str,
    repo: PersonRepositoryInterface = Depends(get_repo),
):
    return repo.list_by_city_state(city, state)

@app.get("/persons/{person_id}", response_model=PersonOut)
def read_person(person_id: str, repo: PersonRepositoryInterface = Depends(get_repo)):
    person = repo.get_person(person_id)
    if not person:
        raise HTTPException(status_code=404, detail="Person not found")
    return person

@app.put("/persons/{person_id}", response_model=PersonOut)
def update_person(
    person_id: str,
    person: PersonUpdate,
    repo: PersonRepositoryInterface = Depends(get_repo),
):
    updated = repo.update_person(person_id, person)
    if not updated:
        raise HTTPException(status_code=404, detail="Person not found")
    return updated

@app.delete("/persons/{person_id}")
def delete_person(person_id: str, repo: PersonRepositoryInterface = Depends(get_repo)):
    success = repo.delete_person(person_id)
    if not success:
        raise HTTPException(status_code=404, detail="Person not found")
    return {"ok": True}
