import os
import pytest

def pytest_addoption(parser):
    parser.addoption(
        "--backend",
        action="store",
        default="sqlite",
        help="Backend to use for tests (sqlite|mongo)",
    )

import shutil

@pytest.fixture(scope="session", autouse=True)
def set_backend_env(request):
    backend = request.config.getoption("--backend")
    os.environ["PERSON_REPO_BACKEND"] = backend
    # Set test DB names/paths
    if backend == "mongo":
        os.environ["MONGO_DB"] = "test_person_db"
        print("[pytest] Using Mongo test DB: test_person_db")
    elif backend == "sqlite":
        db_dir = os.path.dirname("db/test_persons.db")
        if not os.path.exists(db_dir):
            os.makedirs(db_dir)
        os.environ["SQLITE_DB_PATH"] = "db/test_persons.db"
        print("[pytest] Using SQLite test DB: db/test_persons.db")
        print(f"[pytest] Absolute test DB path: {os.path.abspath('db/test_persons.db')}")
    print(f"\n[pytest] Running tests with PERSON_REPO_BACKEND={backend}\n")

    def cleanup():
        if backend == "mongo":
            from pymongo import MongoClient
            uri = os.getenv("MONGO_URI", "mongodb://superuser:supersecret@localhost:27080/person_db?authSource=admin")
            client = MongoClient(uri, serverSelectionTimeoutMS=3000)
            client.drop_database("test_person_db")
            print("[pytest] Dropped Mongo test DB: test_person_db")
        elif backend == "sqlite":
            db_path = "db/test_persons.db"
            # Disabled cleanup for inspection purposes:
            try:
                shutil.move(db_path, db_path+".bak")
                os.remove(db_path+".bak")
                print(f"[pytest] Removed SQLite test DB: {db_path}")
            except FileNotFoundError:
                pass
            print(f"[pytest] (Cleanup skipped) SQLite test DB remains: {db_path}")
    request.addfinalizer(cleanup)
