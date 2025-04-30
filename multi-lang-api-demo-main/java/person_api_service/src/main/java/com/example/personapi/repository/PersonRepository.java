package com.example.personapi.repository;

import com.example.personapi.model.Person;
import java.util.List;
import java.util.Optional;

public interface PersonRepository {
    Person save(Person person);
    Optional<Person> findById(String id);
    List<Person> findByFirstNameAndLastName(String firstName, String lastName);
    List<Person> findByCityAndState(String city, String state);
    List<Person> findAll();
    void deleteById(String id);
}
