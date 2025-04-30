import os

import pytest
from fastapi.testclient import TestClient

from app.main import app

# Use a separate test DB
TEST_DB_PATH = os.path.join(os.path.dirname(__file__), '../db/test_persons.db')

@pytest.fixture(scope="session", autouse=True)
def setup_test_db():
    if os.path.exists(TEST_DB_PATH):
        os.remove(TEST_DB_PATH)
    from app.database import Base, engine
    Base.metadata.create_all(bind=engine)
    yield
    if os.path.exists(TEST_DB_PATH):
        os.remove(TEST_DB_PATH)

@pytest.fixture(scope="module")
def client():
    with TestClient(app) as c:
        yield c

def test_create_and_get_person(client):
    payload = {
        "first_name": "John",
        "last_name": "Doe",
        "age": 30,
        "address": {
            "address_line1": "123 Main St",
            "address_line2": "Apt 1",
            "city": "Springfield",
            "state": "IL",
            "zip": "62704"
        }
    }
    resp = client.post("/persons/", json=payload)
    assert resp.status_code == 200
    data = resp.json()
    assert data["first_name"] == "John"
    person_id = data["id"]
    resp = client.get(f"/persons/{person_id}")
    assert resp.status_code == 200
    assert resp.json()["last_name"] == "Doe"

def test_update_person(client):
    payload = {
        "first_name": "Jane",
        "last_name": "Smith",
        "age": 25,
        "address": {
            "address_line1": "456 Elm St",
            "address_line2": "",
            "city": "Springfield",
            "state": "IL",
            "zip": "62704"
        }
    }
    resp = client.post("/persons/", json=payload)
    person_id = resp.json()["id"]
    update_payload = payload.copy()
    update_payload["age"] = 26
    update_payload["address"]["city"] = "Chicago"
    resp = client.put(f"/persons/{person_id}", json=update_payload)
    assert resp.status_code == 200
    assert resp.json()["age"] == 26
    assert resp.json()["address"]["city"] == "Chicago"

def test_delete_person(client):
    payload = {
        "first_name": "Alice",
        "last_name": "Wonder",
        "age": 22,
        "address": {
            "address_line1": "789 Oak St",
            "address_line2": None,
            "city": "Peoria",
            "state": "IL",
            "zip": "61602"
        }
    }
    resp = client.post("/persons/", json=payload)
    person_id = resp.json()["id"]
    resp = client.delete(f"/persons/{person_id}")
    assert resp.status_code == 200
    resp = client.get(f"/persons/{person_id}")
    assert resp.status_code == 404

def test_search_by_name(client):
    payload = {
        "first_name": "Bob",
        "last_name": "Builder",
        "age": 40,
        "address": {
            "address_line1": "321 Build St",
            "address_line2": None,
            "city": "Metropolis",
            "state": "NY",
            "zip": "10001"
        }
    }
    client.post("/persons/", json=payload)
    resp = client.get("/persons/search?first_name=Bob")
    assert resp.status_code == 200
    results = resp.json()
    assert any(p["last_name"] == "Builder" for p in results)

def test_list_by_city_state(client):
    payload = {
        "first_name": "Clark",
        "last_name": "Kent",
        "age": 35,
        "address": {
            "address_line1": "1 Superman Plaza",
            "address_line2": None,
            "city": "Metropolis",
            "state": "NY",
            "zip": "10001"
        }
    }
    client.post("/persons/", json=payload)
    resp = client.get("/persons/by_city_state?city=Metropolis&state=NY")
    assert resp.status_code == 200
    results = resp.json()
    assert any(p["first_name"] == "Clark" for p in results)
