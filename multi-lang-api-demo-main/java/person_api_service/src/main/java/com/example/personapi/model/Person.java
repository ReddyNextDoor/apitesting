package com.example.personapi.model;

import jakarta.persistence.*;
import lombok.Data;

@Entity
public class Person {
    @Id
    private String id; // Use String for consistency with Python/ObjectId

    private String firstName;
    private String lastName;
    private int age;

    @Embedded
    private Address address;

    public Person() {}

    public String getId() {
        return id;
    }
    public void setId(String id) {
        this.id = id;
    }
    public String getFirstName() {
        return firstName;
    }
    public void setFirstName(String firstName) {
        this.firstName = firstName;
    }
    public String getLastName() {
        return lastName;
    }
    public void setLastName(String lastName) {
        this.lastName = lastName;
    }
    public int getAge() {
        return age;
    }
    public void setAge(int age) {
        this.age = age;
    }
    public Address getAddress() {
        return address;
    }
    public void setAddress(Address address) {
        this.address = address;
    }
}

