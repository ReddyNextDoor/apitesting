import os

from sqlalchemy import create_engine
from sqlalchemy.orm import declarative_base, sessionmaker

# Use SQLITE_DB_PATH env var if set, else default to db/persons.db
DB_PATH = os.getenv('SQLITE_DB_PATH')
if DB_PATH is None:
    DB_FOLDER = os.path.join(os.path.dirname(os.path.dirname(__file__)), 'db')
    if not os.path.exists(DB_FOLDER):
        os.makedirs(DB_FOLDER)
    DB_PATH = os.path.join(DB_FOLDER, 'persons.db')
DB_URL = f'sqlite:///{DB_PATH}'

engine = create_engine(DB_URL, connect_args={"check_same_thread": False})
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)
Base = declarative_base()
