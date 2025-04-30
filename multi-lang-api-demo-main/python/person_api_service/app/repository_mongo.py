import os
from typing import Optional

from bson.objectid import ObjectId
from bson.errors import InvalidId
from loguru import logger
from pymongo import MongoClient
from fastapi import HTTPException

from app.models import PersonCreate, PersonUpdate
from app.repository_interface import PersonRepositoryInterface


class MongoPersonRepository(PersonRepositoryInterface):
    def __init__(self, uri=None, db_name=None):
        db_name = db_name or os.getenv("MONGO_DB", "person_db")
        uri = uri or os.getenv("MONGO_URI", "mongodb://superuser:supersecret@localhost:27080/person_db?authSource=admin")
        try:
            # Mask password for logging
            safe_uri = uri.replace(uri.split(':')[2].split('@')[0], '***') if '://' in uri and '@' in uri else uri
            logger.info(f"Connecting to MongoDB at {safe_uri}")
            self.client = MongoClient(uri, serverSelectionTimeoutMS=3000)
            # Force connection on a request as the connect=True parameter of MongoClient seems
            # to be useless here
            self.client.admin.command('ping')
            logger.info("MongoDB connection established successfully.")
            self.db = self.client[db_name]
            self.collection = self.db["persons"]
        except Exception as e:
            logger.error(f"Failed to connect to MongoDB: {e}")
            raise

    def create_person(self, person: PersonCreate):
        logger.info(f"Input to create_person: {person}")
        try:
            doc = person.model_dump()
            doc["address"] = person.address.model_dump()
            result = self.collection.insert_one(doc)
            doc["id"] = str(result.inserted_id)
            logger.info(f"Output from create_person: {doc}")
            return doc
        except Exception as e:
            logger.error(f"Exception in create_person: {e}")
            raise

    def get_person(self, person_id: str):
        logger.info(f"Input to get_person: {person_id}")
        try:
            try:
                obj_id = ObjectId(person_id)
            except InvalidId:
                logger.error(f"Invalid ObjectId: {person_id}")
                raise HTTPException(status_code=400, detail=f"'{person_id}' is not a valid ID.")
            doc = self.collection.find_one({"_id": obj_id})
            if doc:
                doc["id"] = str(doc["_id"])
                logger.info(f"Output from get_person: {doc}")
                return doc
            logger.info(f"Person not found: {person_id}")
            raise HTTPException(status_code=404, detail=f"Person with id '{person_id}' not found.")
        except HTTPException:
            raise
        except Exception as e:
            logger.error(f"Exception in get_person: {e}")
            raise HTTPException(status_code=500, detail="Internal server error.")

    def update_person(self, person_id: str, person: PersonUpdate):
        logger.info(f"Input to update_person: {person_id}, {person}")
        try:
            try:
                obj_id = ObjectId(person_id)
            except InvalidId:
                logger.error(f"Invalid ObjectId: {person_id}")
                raise HTTPException(status_code=400, detail=f"'{person_id}' is not a valid ID.")
            update_doc = person.model_dump()
            update_doc["address"] = person.address.model_dump()
            result = self.collection.update_one({"_id": obj_id}, {"$set": update_doc})
            if result.matched_count:
                doc = self.collection.find_one({"_id": obj_id})
                doc["id"] = str(doc["_id"])
                logger.info(f"Output from update_person: {doc}")
                return doc
            logger.info(f"Person not found for update: {person_id}")
            raise HTTPException(status_code=404, detail=f"Person with id '{person_id}' not found.")
        except HTTPException:
            raise
        except Exception as e:
            logger.error(f"Exception in update_person: {e}")
            raise HTTPException(status_code=500, detail="Internal server error.")

    def delete_person(self, person_id: str) -> bool:
        logger.info(f"Input to delete_person: {person_id}")
        try:
            try:
                obj_id = ObjectId(person_id)
            except InvalidId:
                logger.error(f"Invalid ObjectId: {person_id}")
                raise HTTPException(status_code=400, detail=f"'{person_id}' is not a valid ID.")
            result = self.collection.delete_one({"_id": obj_id})
            if result.deleted_count > 0:
                logger.info(f"Output from delete_person: True")
                return True
            logger.info(f"Person not found for delete: {person_id}")
            raise HTTPException(status_code=404, detail=f"Person with id '{person_id}' not found.")
        except HTTPException:
            raise
        except Exception as e:
            logger.error(f"Exception in delete_person: {e}")
            raise HTTPException(status_code=500, detail="Internal server error.")

    def search_by_name(self, first_name: Optional[str], last_name: Optional[str]):
        logger.info(f"Input to search_by_name: first_name={first_name}, last_name={last_name}")
        try:
            query = {}
            if first_name:
                query["first_name"] = {"$regex": first_name, "$options": "i"}
            if last_name:
                query["last_name"] = {"$regex": last_name, "$options": "i"}
            results = list(self.collection.find(query))
            for doc in results:
                doc["id"] = str(doc["_id"])
            logger.info(f"Output from search_by_name: {results}")
            return results
        except Exception as e:
            logger.error(f"Exception in search_by_name: {e}")
            raise

    def list_by_city_state(self, city: str, state: str):
        logger.info(f"Input to list_by_city_state: city={city}, state={state}")
        try:
            query = {"address.city": city, "address.state": state}
            results = list(self.collection.find(query))
            for doc in results:
                doc["id"] = str(doc["_id"])
            logger.info(f"Output from list_by_city_state: {results}")
            return results
        except Exception as e:
            logger.error(f"Exception in list_by_city_state: {e}")
            raise
