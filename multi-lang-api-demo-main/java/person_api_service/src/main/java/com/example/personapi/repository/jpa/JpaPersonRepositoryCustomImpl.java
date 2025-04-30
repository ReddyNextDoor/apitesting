package com.example.personapi.repository.jpa;

import com.example.personapi.model.Person;
import com.example.personapi.repository.PersonRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Profile;
import org.springframework.stereotype.Component;
import java.util.List;
import java.util.Optional;

@Component
@Profile("sqlite")
public class JpaPersonRepositoryCustomImpl implements PersonRepository {
    private final JpaPersonRepository jpaPersonRepository;

    @Autowired
    public JpaPersonRepositoryCustomImpl(JpaPersonRepository jpaPersonRepository) {
        this.jpaPersonRepository = jpaPersonRepository;
    }

    @Override
    public Person save(Person person) {
        return jpaPersonRepository.save(person);
    }

    @Override
    public Optional<Person> findById(String id) {
        return jpaPersonRepository.findById(id);
    }

    @Override
    public List<Person> findByFirstNameAndLastName(String firstName, String lastName) {
        return jpaPersonRepository.findByFirstNameContainingIgnoreCaseAndLastNameContainingIgnoreCase(firstName, lastName);
    }

    @Override
    public List<Person> findByCityAndState(String city, String state) {
        return jpaPersonRepository.findByAddress_CityAndAddress_State(city, state);
    }

    @Override
    public List<Person> findAll() {
        return jpaPersonRepository.findAll();
    }

    @Override
    public void deleteById(String id) {
        jpaPersonRepository.deleteById(id);
    }
}
