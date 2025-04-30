package com.example.personapi.repository.jpa;

import com.example.personapi.model.Person;
import org.springframework.data.jpa.repository.JpaRepository;
import java.util.List;

public interface JpaPersonRepository extends JpaRepository<Person, String> {
    List<Person> findByFirstNameContainingIgnoreCaseAndLastNameContainingIgnoreCase(String firstName, String lastName);
    List<Person> findByAddress_CityAndAddress_State(String city, String state);
}
