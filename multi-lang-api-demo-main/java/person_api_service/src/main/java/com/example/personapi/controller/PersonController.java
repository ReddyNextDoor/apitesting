package com.example.personapi.controller;

import java.util.UUID;

import com.example.personapi.dto.PersonDto;
import com.example.personapi.model.Person;
import com.example.personapi.repository.PersonRepositoryFactory;
import com.example.personapi.repository.PersonRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Optional;
import java.util.stream.Collectors;

@RestController
@RequestMapping("/persons")
public class PersonController {
    private final PersonRepositoryFactory factory;
    private PersonRepository personRepository;

    @Autowired
    public PersonController(PersonRepositoryFactory factory) {
        this.factory = factory;
    }


    @GetMapping
    public List<PersonDto> getAllPersons() {
        return getPersonRepository().findAll().stream().map(this::toDto).collect(Collectors.toList());
    }

    @GetMapping("/{id}")
    public PersonDto getPersonById(@PathVariable String id) {
        Optional<Person> person = getPersonRepository().findById(id);
        return person.map(this::toDto).orElse(null);
    }

    @PostMapping
    public PersonDto createPerson(@RequestBody PersonDto personDto) {
        Person person = toEntity(personDto);
        if (person.getId() == null || person.getId().isEmpty()) {
            person.setId(UUID.randomUUID().toString());
        }
        Person savedPerson = getPersonRepository().save(person);
        return toDto(savedPerson);
    }

    @PutMapping("/{id}")
    public PersonDto updatePerson(@PathVariable String id, @RequestBody PersonDto personDto) {
        Person person = toEntity(personDto);
        person.setId(id);
        Person updated = getPersonRepository().save(person);
        return toDto(updated);
    }

    @DeleteMapping("/{id}")
    public void deletePerson(@PathVariable String id) {
        getPersonRepository().deleteById(id);
    }

    @GetMapping("/search")
    public List<PersonDto> searchPersons(@RequestParam String firstName, @RequestParam String lastName) {
        return getPersonRepository().findByFirstNameAndLastName(firstName, lastName)
                .stream().map(this::toDto).collect(Collectors.toList());
    }

    @GetMapping("/by_city_state")
    public List<PersonDto> getByCityState(@RequestParam String city, @RequestParam String state) {
        return getPersonRepository().findByCityAndState(city, state)
                .stream().map(this::toDto).collect(Collectors.toList());
    }

    // --- Mapping helpers ---
    private PersonDto toDto(Person person) {
        PersonDto dto = new PersonDto();
        if (person.getId() != null) {
            dto.setId(person.getId().toString());
        } else {
            dto.setId(null);
        }
        dto.setFirstName(person.getFirstName());
        dto.setLastName(person.getLastName());
        dto.setAge(person.getAge());
        dto.setAddress(person.getAddress());
        return dto;
    }

    private Person toEntity(PersonDto dto) {
        Person p = new Person();
        if (dto.getId() != null && !dto.getId().isEmpty()) {
            p.setId(dto.getId());
        }
        p.setFirstName(dto.getFirstName());
        p.setLastName(dto.getLastName());
        p.setAge(dto.getAge());
        p.setAddress(dto.getAddress());
        return p;
    }

    private PersonRepository getPersonRepository() {
        if (personRepository == null) {
            personRepository = factory.getPersonRepository();
        }
        return personRepository;
    }
}

