from typing import List, Optional

from loguru import logger
from sqlalchemy.orm import Session

from app.models import AddressORM, PersonCreate, PersonORM, PersonUpdate, PersonOut
from sqlalchemy import and_
from app.repository_interface import PersonRepositoryInterface


def orm_to_personout_dict(person_orm):
    address_orm = person_orm.address
    address_dict = None
    if address_orm:
        address_dict = {
            "address_line1": address_orm.address_line1,
            "address_line2": address_orm.address_line2,
            "city": address_orm.city,
            "state": address_orm.state,
            "zip": address_orm.zip,
        }
    person_dict = {
        "id": str(person_orm.id),
        "first_name": person_orm.first_name,
        "last_name": person_orm.last_name,
        "age": person_orm.age,
        "address": address_dict,
    }
    return person_dict

import os

class PersonRepository(PersonRepositoryInterface):
    def __init__(self, db: Session = None):
        # If db is not provided, create a new engine/session using the env var
        if db is not None:
            self.db = db
        else:
            from sqlalchemy import create_engine
            from sqlalchemy.orm import sessionmaker
            db_path = os.getenv("SQLITE_DB_PATH", "db/persons.db")
            engine = create_engine(f"sqlite:///{db_path}")
            SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)
            self.db = SessionLocal()

    def create_person(self, person: PersonCreate):
        logger.info(f"Input to create_person: {person}")
        address_data = person.address.model_dump()
        db_address = AddressORM(**address_data)
        db_person = PersonORM(
            first_name=person.first_name,
            last_name=person.last_name,
            age=person.age,
            address=db_address
        )
        self.db.add(db_person)
        self.db.commit()
        self.db.refresh(db_person)
        result = PersonOut.model_validate(orm_to_personout_dict(db_person)).model_dump()
        logger.info(f"Output from create_person: {result}")
        return result

    def get_person(self, person_id: int):
        logger.info(f"Input to get_person: {person_id}")
        person = self.db.query(PersonORM).filter(PersonORM.id == person_id).first()
        if person:
            result = PersonOut.model_validate(orm_to_personout_dict(person)).model_dump()
            logger.info(f"Output from get_person: {result}")
            return result
        logger.info(f"Output from get_person: None")
        return None

    def update_person(self, person_id: int, person: PersonUpdate):
        logger.info(f"Input to update_person: {person_id}, {person}")
        db_person = self.db.query(PersonORM).filter(PersonORM.id == person_id).first()
        if not db_person:
            logger.info("Output from update_person: None (not found)")
            return None
        person_data = person.model_dump()
        for key, value in person_data.items():
            if key == "address" and value is not None:
                if db_person.address:
                    for addr_key, addr_value in value.items():
                        setattr(db_person.address, addr_key, addr_value)
                else:
                    db_person.address = AddressORM(**value)
            else:
                setattr(db_person, key, value)
        self.db.commit()
        self.db.refresh(db_person)
        result = PersonOut.model_validate(orm_to_personout_dict(db_person)).model_dump()
        logger.info(f"Output from update_person: {result}")
        return result

    def delete_person(self, person_id: int) -> bool:
        logger.info(f"Input to delete_person: {person_id}")
        db_person = self.db.query(PersonORM).filter(PersonORM.id == person_id).first()
        if not db_person:
            logger.info("Output from delete_person: False (not found)")
            return False
        self.db.delete(db_person)
        self.db.commit()
        logger.info("Output from delete_person: True")
        return True

    def search_by_name(self, first_name: Optional[str], last_name: Optional[str]):
        logger.info(f"Input to search_by_name: first_name={first_name}, last_name={last_name}")
        query = self.db.query(PersonORM)
        if first_name:
            query = query.filter(PersonORM.first_name.ilike(f"%{first_name}%"))
        if last_name:
            query = query.filter(PersonORM.last_name.ilike(f"%{last_name}%"))
        persons = query.all()
        results = [PersonOut.model_validate(orm_to_personout_dict(person)).model_dump() for person in persons]
        logger.info(f"Output from search_by_name: {results}")
        return results

    def list_by_city_state(self, city: str, state: str):
        logger.info(f"Input to list_by_city_state: city={city}, state={state}")
        persons = self.db.query(PersonORM).filter(
            PersonORM.address.has(
                and_(
                    AddressORM.city == city,
                    AddressORM.state == state
                )
            )
        ).all()
        results = [PersonOut.model_validate(orm_to_personout_dict(person)).model_dump() for person in persons]
        logger.info(f"Output from list_by_city_state: {results}")
        return results
