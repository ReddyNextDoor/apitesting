package com.example.personapi.repository.mongo;

import com.example.personapi.model.Person;
import com.example.personapi.repository.PersonRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Profile;
import java.util.List;
import java.util.Optional;
import org.springframework.stereotype.Component;

@Component
@Profile("mongo")
public class MongoPersonRepositoryCustomImpl implements PersonRepository {
    private final MongoPersonRepository mongoPersonRepository;

    @Autowired
    public MongoPersonRepositoryCustomImpl(MongoPersonRepository mongoPersonRepository) {
        this.mongoPersonRepository = mongoPersonRepository;
    }

    @Override
    public Person save(Person person) {
        return mongoPersonRepository.save(person);
    }

    @Override
    public Optional<Person> findById(String id) {
        return mongoPersonRepository.findById(id);
    }

    @Override
    public List<Person> findByFirstNameAndLastName(String firstName, String lastName) {
        return mongoPersonRepository.findByFirstNameContainingIgnoreCaseAndLastNameContainingIgnoreCase(firstName, lastName);
    }

    @Override
    public List<Person> findByCityAndState(String city, String state) {
        return mongoPersonRepository.findByAddress_CityAndAddress_State(city, state);
    }

    @Override
    public List<Person> findAll() {
        return mongoPersonRepository.findAll();
    }

    @Override
    public void deleteById(String id) {
        mongoPersonRepository.deleteById(id);
    }
}
