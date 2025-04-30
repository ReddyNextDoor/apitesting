package com.example.personapi.repository.mongo;

import com.example.personapi.model.Person;
import org.springframework.data.mongodb.repository.MongoRepository;
import java.util.List;

public interface MongoPersonRepository extends MongoRepository<Person, String> {
    List<Person> findByFirstNameContainingIgnoreCaseAndLastNameContainingIgnoreCase(String firstName, String lastName);
    List<Person> findByAddress_CityAndAddress_State(String city, String state);
}
