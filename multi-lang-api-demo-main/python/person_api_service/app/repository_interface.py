from abc import ABC, abstractmethod
from typing import Optional

from app.models import PersonCreate, PersonUpdate


class PersonRepositoryInterface(ABC):
    @abstractmethod
    def create_person(self, person: PersonCreate):
        pass

    @abstractmethod
    def get_person(self, person_id: int):
        pass

    @abstractmethod
    def update_person(self, person_id: int, person: PersonUpdate):
        pass

    @abstractmethod
    def delete_person(self, person_id: int) -> bool:
        pass

    @abstractmethod
    def search_by_name(self, first_name: Optional[str], last_name: Optional[str]):
        pass

    @abstractmethod
    def list_by_city_state(self, city: str, state: str):
        pass
