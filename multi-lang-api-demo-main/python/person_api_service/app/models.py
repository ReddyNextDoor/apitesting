from typing import Optional

from pydantic import BaseModel
from sqlalchemy import Column, ForeignKey, Integer, String
from sqlalchemy.orm import relationship

from app.database import Base


class AddressORM(Base):
    __tablename__ = 'addresses'
    id = Column(Integer, primary_key=True, index=True)
    address_line1 = Column(String, nullable=False)
    address_line2 = Column(String, nullable=True)
    city = Column(String, nullable=False)
    state = Column(String, nullable=False)
    zip = Column(String, nullable=False)
    person_id = Column(Integer, ForeignKey('persons.id'))

class PersonORM(Base):
    __tablename__ = 'persons'
    id = Column(Integer, primary_key=True, index=True)
    first_name = Column(String, index=True)
    last_name = Column(String, index=True)
    age = Column(Integer)
    address = relationship(
        'AddressORM',
        uselist=False,
        backref='person',
        cascade="all, delete-orphan",
    )

# Pydantic Schemas
class Address(BaseModel):
    address_line1: str
    address_line2: Optional[str] = None
    city: str
    state: str
    zip: str

class PersonBase(BaseModel):
    first_name: str
    last_name: str
    age: int
    address: Address

class PersonCreate(PersonBase):
    pass

class PersonUpdate(PersonBase):
    pass

class PersonOut(PersonBase):
    id: str
    model_config = {"from_attributes": True}
